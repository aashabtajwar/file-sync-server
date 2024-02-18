eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFhc2hhYkBnbWFpbC5jb20iLCJpZCI6IjEiLCJleHAiOjE3MDgxODIwMTl9.FFQeBBSdIyFlmhZKTGVHxZJMftNr_rHIax9vNEcNABodef comb_two(nums, target):
    res = []
    
    def backtrack(num, curr):
        # print(f'Current Num Value = {num}')
        # base cases
        if sum(curr) > target:
            print(f'return Current Array = {curr}')
            return  
        if sum(curr) == target:
            x = curr.copy()
            x.sort()
            if x not in res:
                res.append(x)
            print(f'Update = {res}')
            return 

        if num == len(nums):
            if sum(curr) == target:
                x = curr.sort()
                x.sort()
                if x not in res:
                    res.append(x)
                    print(f'Update = {res}')
            return
        
        for i in range(num, len(nums)):
            print(f'Current Array => {curr}, i = {i}')
            curr.append(nums[i])
            backtrack(i + 1, curr)
            curr.pop()    
            
    backtrack(0, [])
    return res

nums = [32,10,32,5,25,9,18,23,28,24,10,33,6,24,32,18,10,28,17,18,13,22,7,25,22,17,28,13,17,32,19,6,7,17,7,28,21,12,8,6,31,13,34,24,24,31,8,29,16,20,12,25,29,6,15,16,19,30,17,23,27,31,17,29]
target = 30
print(comb_two(nums, target))