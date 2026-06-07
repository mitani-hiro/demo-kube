package trino.access_control

import rego.v1

# デフォルトでは拒否
default allow := false

# ユーザーとロールの定義
user_roles := {
	"alice": ["admin", "data_analyst"],
	"bob": ["data_analyst"],
	"charlie": ["data_viewer"],
	"david": ["finance_user"],
}

# Trino の OPA プラグインが送る実際の operation 名で定義する
# （SELECT/INSERT といった SQL 名ではなく ExecuteQuery 等が渡ってくる）
read_operations := {
	"ExecuteQuery",
	"AccessCatalog",
	"SelectFromColumns",
	"FilterCatalogs",
	"FilterSchemas",
	"FilterTables",
	"FilterColumns",
	"ShowSchemas",
	"ShowTables",
	"ShowColumns",
	"GetColumnMask",
	"GetRowFilters",
	"ExecuteFunction",
	"FilterFunctions",
}

write_operations := {
	"InsertIntoTable",
	"DeleteFromTable",
	"UpdateTableColumns",
	"TruncateTable",
}

# admin は全操作を許可
operation_allowed("admin", _)

operation_allowed("data_analyst", operation) if {
	operation in (read_operations | write_operations)
}

operation_allowed("data_viewer", operation) if {
	operation in read_operations
}

operation_allowed("finance_user", operation) if {
	operation in read_operations
}

allow if {
	user := input.context.identity.user
	some role in user_roles[user]
	operation_allowed(role, input.action.operation)
}
