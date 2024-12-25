import random

def generate_nodes(filename, count=500, max_id=1024):
    nodes = random.sample(range(max_id), count)
    with open(filename, "w") as file:
        for node in sorted(nodes):
            file.write(f"{node}\n")

if __name__ == "__main__":
    generate_nodes("nodes.txt")