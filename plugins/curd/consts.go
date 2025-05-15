package curd

type JoinType string

// JoinLeft ...联合查询相关
const (
	JoinLeft  JoinType = "left"
	JoinRight JoinType = "right"
	JoinInner JoinType = "inner"
)
