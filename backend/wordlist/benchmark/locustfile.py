# -*- coding:utf-8 -*- #

# Usage:
#  locust --host=http://localhost:8006

from locust import HttpLocust, TaskSet, task

class UserBehavior(TaskSet):
  @task(1)
  def predict(self):
    text = "哈哈"
    self.client.post("/classify", {"text":text})

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait=2000
    max_wait=2000