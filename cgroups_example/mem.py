import time

buffer = []

def main():
    while True:
        buffer.append(" " * 100 * 1024 * 1024)
        print(f"[{time.time()}] :: 100 MiB was allocated")

if __name__ == "__main__":
    main()