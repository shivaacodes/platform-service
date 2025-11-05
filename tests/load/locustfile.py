from locust import HttpUser, between, task

class PlatformServiceUser(HttpUser):
    wait_time = between(0.2, 0.5) # realistic pacing

    @task(3)
    def get_data(self):
        self.client.get("/api/v1/data")

    @task(1)
    def check_ready(self):
        self.client.get("/readyz")

