eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFhc2hhYkBnbWFpbC5jb20iLCJpZCI6IjEiLCJleHAiOjE3MDgyMzU1NDZ9.XyfxfHK3r4FrA2e9rd6JdRHhKPKW-erscW-ryhe9OsA# Source to Target

def source_target(graph):
    res = []

    def dfs(num, curr):
        if num == len(graph) - 1:
            res.append(curr.copy())
        for i in graph[num]:
            curr.append(i)
            dfs(i, curr)
            curr.pop()
    dfs(0, [])
    for g in res:
        g.insert(0, 0)
    return res


graph = [[4,3,1],[3,2,4],[3],[4],[]]
print(source_target(graph))