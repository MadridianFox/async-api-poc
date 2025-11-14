from locust import LoadTestShape

class StepLoadShape(LoadTestShape):
    abstract = True

    def __init__(self):
        super().__init__()
        # (duration, users)
        self.stages: list[tuple[int,int]]
        self.total_duration = sum(d for d, _ in self.stages)
        self.spawn_rate = 1.0


    def tick(self):
        run_time = self.get_run_time()

        elapsed = 0
        prev_users = 0
        for duration, target_users in self.stages:
            start_time = elapsed
            end_time = elapsed + duration
            if run_time < end_time:
                progress = (run_time - start_time) / duration
                users = prev_users + (target_users - prev_users) * progress
                return round(users), self.spawn_rate
            prev_users = target_users
            elapsed = end_time

        return None