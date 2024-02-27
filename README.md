# sql2pb

> Generates a protobuf file from your database.

## Uses

```shell
Generates a protobuf file from your database

Usage:
  sql2pb gen [flags]

Flags:
      --db_type string           the database type. mysql | postgres (default "mysql")
      --dbname string            the database name
      --field_style string       gen protobuf field style. sql_pb | sqlPb (default "sql_pb")
      --go_package string        the protocol buffer go_package. defaults to the database schema.
  -h, --help                     help for gen
      --host string              the database host (default "localhost")
      --ignore_columns strings   a comma spaced list of mysql columns to ignore
      --ignore_tables strings    a comma spaced list of tables to ignore
      --package string           the protocol buffer package. defaults to the database schema.
      --password string          the database password
      --port int                 the database port (default 3306)
      --schema string            the database schema
      --service_name string      the protocol buffer package. defaults to the database schema.
      --table string             the table schema. multiple tables ',' split. 
      --user string              the database user (default "root")

```

```shell
sql2pb gen  --host=127.0.0.1 --port=3306 --dbname=root --user=root --password=123456  --service_name=User --db_type=mysql --table=sys_user --go_package=./pb --package=user
```

```protobuf
syntax = "proto3";

option go_package = "./pb";

package user;

// ------------------------------------
// Messages
// ------------------------------------

//--------------------------------用户--------------------------------

message SysUser {
    int64 id = 1; // ID
    string username = 2; // 用户名
    string password = 3; // 密码
    int64 create_at = 4; // 创建时间
    int64 update_at = 5; // 修改时间
    int64 delete_at = 6; // 删除时间
}

message SysUserFilter {
    optional int64 id = 1; // ID
    optional string username = 2; // 用户名
    optional string password = 3; // 密码
    optional int64 create_at = 4; // 创建时间
    optional int64 update_at = 5; // 修改时间
    optional int64 delete_at = 6; // 删除时间
}

message AddSysUserReq {
    string username = 1; // 用户名
    string password = 2; // 密码
}

message AddSysUserResp {
}

message UpdateSysUserReq {
    optional int64 id = 1; // ID
    optional string username = 2; // 用户名
    optional string password = 3; // 密码
}

message UpdateSysUserResp {
}

message DelSysUserReq {
    int64 id = 1; // id
}

message DelSysUserResp {
}

message SelectSysUserByIdReq {
    int64 id = 1; // id
}

message SelectSysUserByIdResp {
    SysUser sys_user = 1; // sys_user
}

message SelectSysUserListReq {
    int64 page = 1; // 页码
    int64 page_size = 2; // 每页数量
    optional SysUserFilter filter = 3; // SysUserFilter
}

message SelectSysUserListResp {
    int64 count = 1; // 总数
    int64 page_count = 2; // 页码总数
    repeated SysUser results = 3; // sys_user
}

// ------------------------------------
// Rpc Func
// ------------------------------------

service User {

    //-----------------------用户-----------------------

    // 创建用户
    rpc InsertSysUser (AddSysUserReq) returns (AddSysUserResp);

    // 更新用户
    rpc UpdateSysUser (UpdateSysUserReq) returns (UpdateSysUserResp);

    // 根据 用户 id 删除
    rpc DeleteSysUser (DelSysUserReq) returns (DelSysUserResp);

    // 根据 用户 id 获取详情
    rpc SelectSysUserById (SelectSysUserByIdReq) returns (SelectSysUserByIdResp);

    // 用户 列表
    rpc SelectSysUserList (SelectSysUserListReq) returns (SelectSysUserListResp);
}

```

## Thanks

[https://github.com/Mikaelemmmm/sql2pb](https://github.com/Mikaelemmmm/sql2pb)