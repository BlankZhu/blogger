/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	} else if l2 == nil {
		return l1
	}

	// pad to the same length
	p1, p2 := l1, l2
	for {
		if p1.Next == nil && p2.Next == nil {
			break
		} else if p1.Next == nil && p2.Next != nil {
			p1.Next = new(ListNode)
			p1.Next.Val, p1.Next.Next = 0, nil
		} else if p1.Next != nil && p2.Next == nil {
			p2.Next = new(ListNode)
			p2.Next.Val, p2.Next.Next = 0, nil
		}
		p1, p2 = p1.Next, p2.Next
	}

	p1, p2 = l1, l2
	head := new(ListNode)
	head.Next = nil
	head.Val = 0
	forward := 0
	for p1 != nil && p2 != nil {
		sum := p1.Val + p2.Val + forward
		if sum >= 10 {
			forward, p1.Val = sum/10, sum%10
		} else {
			p1.Val, forward = sum, 0
		}
		p1, p2 = p1.Next, p2.Next
	}

	if forward != 0 {
		tmp := l1
		for tmp.Next != nil {
			tmp = tmp.Next
		}
		head.Val, tmp.Next = forward, head
	}

	return l1
}