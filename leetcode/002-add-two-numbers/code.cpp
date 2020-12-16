struct ListNode {
    int val;
    ListNode *next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode *next) : val(x), next(next) {}
};

class Solution {
public:
    ListNode* addTwoNumbers(ListNode* l1, ListNode* l2) {
        if (l1 == nullptr) {
            return l2;
        }
        if (l2 == nullptr) {
            return l1;
        }

        auto p1 = l1;
        auto p2 = l2;
        for (;;) {
            if (p1 == nullptr && p2 == nullptr) {
                break;
            } else if (p1->next == nullptr && p2->next != nullptr) {
                p1->next = new ListNode;
            } else if (p1->next != nullptr && p2->next == nullptr) {
                p2->next = new ListNode;
            }
            p1 = p1->next;
            p2 = p2->next;
        }

        p1 = l1;
        p2 = l2;
        ListNode* tmp;
        int forward = 0;
        while (p1 != nullptr && p2 != nullptr) {
            auto sum = p1->val + p2->val + forward;
            if (sum >= 10) {
                p1->val = sum % 10;
                forward = sum / 10;
            } else {
                p1->val = sum;
                forward = 0;
            }
            if(p1->next == nullptr) {
                tmp = p1;
            }
            p1 = p1->next;
            p2 = p2->next;
        }
        if (forward != 0) {
            tmp->next = new ListNode;
            tmp->next->val = forward;
        }

        return l1;
    }
};