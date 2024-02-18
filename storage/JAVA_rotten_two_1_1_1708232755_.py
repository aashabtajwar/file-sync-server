eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFhc2hhYkBnbWFpbC5jb20iLCJpZCI6IjEiLCJleHAiOjE3MDgyMzU1NDZ9.XyfxfHK3r4FrA2e9rd6JdRHhKPKW-erscW-ryhe9OsA
def rotting_oranges(grid):
    ROWS = len(grid)
    COLS = len(grid[0])
    fresh = 0
    minutes = 0
    queue = []
    for r in range(ROWS):
        for c in range(COLS):
            if grid[r][c] == 1:
                fresh += 1
            if grid[r][c] == 2:
                queue.append([r, c])
    
    while len(queue) and fresh > 0:
        for i in range(len(queue)):
            r, c = queue.pop(0)
            neighbours = [[r + 1, c], [r - 1, c], [r, c + 1], [r, c - 1]]
            for n in neighbours:
                if n[0] == ROWS or n[0] < 0 or n[1] == COLS or n[1] < 0 or grid[n[0]][n[1]] != 1:
                    continue
                grid[n[0]][n[1]] = 2
                queue.append(n)
                fresh -= 1
        minutes += 1

    if fresh > 0:
        return -1
    return minutes


grid = [[2,1,1],[1,1,0],[0,1,1]]
print(rotting_oranges(grid))