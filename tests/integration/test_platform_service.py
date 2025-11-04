import time
import requests

BASE_URL = "http://localhost:8080"


def test_readyz():
    """Check Redis connectivity readiness probe."""
    resp = requests.get(f"{BASE_URL}/readyz", timeout=2)
    assert resp.status_code == 200
    assert resp.text.strip().lower() == "ok"


def test_data_endpoint_cache_behavior():
    """Ensure cache-aside pattern works correctly."""
    # First request → MISS expected
    r1 = requests.get(f"{BASE_URL}/api/v1/data", timeout=2)
    assert r1.status_code == 200
    assert r1.headers.get("X-Cache") == "MISS"
    data_1 = r1.json()

    # Second request → HIT expected (cached)
    r2 = requests.get(f"{BASE_URL}/api/v1/data", timeout=2)
    assert r2.status_code == 200
    assert r2.headers.get("X-Cache") == "HIT"
    data_2 = r2.json()

    # Cached value should match previous response
    assert data_1 == data_2

    # Wait for TTL expiry (60 seconds)
    time.sleep(61)
    # Third request → new MISS and new timestamp
    r3 = requests.get(f"{BASE_URL}/api/v1/data", timeout=2)
    assert r3.status_code == 200
    assert r3.headers.get("X-Cache") == "MISS"
    data_3 = r3.json()

    assert data_3["now"] != data_1["now"]

