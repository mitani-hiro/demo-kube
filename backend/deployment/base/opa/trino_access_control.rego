package trino.access_control

import rego.v1

# デフォルトでは拒否
default allow = false

# ユーザーとロールの定義
user_roles := {
    "alice": ["admin", "data_analyst"],
    "bob": ["data_analyst"],
    "charlie": ["data_viewer"],
    "david": ["finance_user"]
}

# テーブルアクセス権限の定義
table_permissions := {
    "admin": {
        "operations": ["SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP"]
    },
    "data_analyst": {
        "operations": ["SELECT", "INSERT", "UPDATE"]
    },
    "data_viewer": {
        "operations": ["SELECT"]
    },
    "finance_user": {
        "operations": ["SELECT"]
    }
}

# ユーザーがロールを持っているかチェック
has_role(user, role) if {
    role in user_roles[user]
}

# ロールが操作を実行できるかチェック
role_can_perform_operation(role, operation) if {
    operation in table_permissions[role].operations
}

# メインのallow判定 - Trinoが期待する形式
allow if {
    print("#### input: ", input)

    user := input.context.identity.user
    operation := input.action.operation

    some role in user_roles[user]
    role_can_perform_operation(role, operation)
}
