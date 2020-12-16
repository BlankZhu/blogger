#include <vector>
#include <unordered_map>

class Solution {
public:
    std::vector<int> twoSum(std::vector<int>& nums, int target) {
        auto map = std::unordered_map<int, int>();
        for(auto i = 0; i < nums.size(); i++) {
            auto val = nums[i];
            auto sub = target - val;
            auto found = map.find(sub);
            if (found != map.end()) {
                auto ret = std::vector<int>(0, 2);
                ret[0] = i;
                ret[1] = val;
                return ret;
            }
            map[val] = i;
        }
        return std::vector<int>(0, 2);
    }
};