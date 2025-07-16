import os
import time

FORK_TIMES: int = 10

def main():
    for _ in range(FORK_TIMES):
        pid = os.fork()
        print(f"[{time.time()}] :: new proprocess with PID: {pid}")

if __name__ == "__main__":
    main()