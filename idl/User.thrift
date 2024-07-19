namespace go user

//table_name:Users
struct User {
    1: i64 Id (idl.tag="primary_key, auto_increment")
    2: string Username (idl.tag="type=varchar(255)")
    3: i32 Age (idl.tag="type=int")
    4: string City (idl.tag="type=varchar(255)") // 城市
    5: string Banned (idl.tag="type=boolean")
    6: i64 RoleId (idl.tag="type=bigint")
    7: string Email (idl.tag="type=varchar(255)")
    8: i32 DefaultAge (idl.tag="type=int, default_value=18")
}
