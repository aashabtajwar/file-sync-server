eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFhc2hhYkBnbWFpbC5jb20iLCJpZCI6IjEiLCJleHAiOjE3MDgyMjc4OTV9.sikVCTP4x43TBnEaTDaC8zpI8XVfVCCoNSDb7z0PnyIdef comb_two(candidates, target):
    res = []
    candidates.sort()
    def backtrack(idx, curr):
        
        if sum(curr) > target: return
        if sum(curr) == target:
            res.append(curr)
            return  
        print(f'Current Curr = {curr}')
        for i in range(idx, len(candidates)):
            if i > idx and candidates[i] == candidates[i - 1]:
                continue
            # curr.append(candidates[i])
            backtrack(i + 1, curr+[candidates[i]])
    
    backtrack(0, [])
    return res


print(comb_two([10,1,2,7,6,1,5], 8))