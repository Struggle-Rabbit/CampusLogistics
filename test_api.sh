#!/bin/bash
# CampusLogistics API 接口测试脚本 v3 - 最终版

BASE_URL="http://127.0.0.1:8880/api/v1"
PASS=0
FAIL=0
TOTAL=0
SKIP=0

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

report() {
    local name="$1" expected="$2" body="$3" status="$4"
    TOTAL=$((TOTAL + 1))
    if echo "$body" | grep -q "$expected"; then
        echo -e "  ${GREEN}[PASS]${NC} $name (HTTP $status)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}[FAIL]${NC} $name (HTTP $status)"
        echo "    期望包含: $expected"
        echo "    实际响应: $body"
        FAIL=$((FAIL + 1))
    fi
}
report_skip() { local name="$1"; TOTAL=$((TOTAL + 1)); SKIP=$((SKIP + 1)); echo -e "  ${BLUE}[SKIP]${NC} $name"; }
check() {
    local name="$1" condition="$2" body="$3" status="$4"
    TOTAL=$((TOTAL + 1))
    if eval "$condition"; then
        echo -e "  ${GREEN}[PASS]${NC} $name (HTTP $status)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}[FAIL]${NC} $name (HTTP $status)"
        echo "    实际响应: $body"
        FAIL=$((FAIL + 1))
    fi
}

echo "========================================"
echo "  CampusLogistics API 接口测试报告"
echo "  时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "  地址: $BASE_URL"
echo "========================================"
echo ""

# ==========================================
# 第一部分：公共接口
# ==========================================
echo -e "${YELLOW}[第一部分] 公共接口（无需认证）${NC}"
echo "---"

echo "  1.1 GET /notice/public"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/notice/public")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "获取公共公告列表" "获取成功" "$BODY" "$STATUS"

echo "  1.2 POST /register"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"测试用户","mobile":"13800138001","password":"Test12345678","userType":"02"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "用户注册(或已存在)" 'echo "$BODY" | grep -qE "注册成功|手机号已注册"' "$BODY" "$STATUS"

echo "  1.3 POST /register (重复)"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{"name":"测试用户","mobile":"13800138001","password":"Test12345678","userType":"02"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "重复注册被拒绝" 'echo "$BODY" | grep -q "手机号已注册"' "$BODY" "$STATUS"

echo "  1.4 POST /login"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"account":"13800138001","password":"Test12345678"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "用户登录" "accessToken" "$BODY" "$STATUS"
ACCESS_TOKEN=$(echo "$BODY" | grep -o '"accessToken":"[^"]*"' | cut -d'"' -f4)
echo "    Token: ${ACCESS_TOKEN:0:40}..."

echo "  1.5 POST /login (密码错误)"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"account":"13800138001","password":"wrongpassword"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "密码错误被拒绝" 'echo "$BODY" | grep -q "账号密码不正确"' "$BODY" "$STATUS"

echo "  1.6 POST /login (账号不存在)"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"account":"99999999999","password":"Test12345678"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "账号不存在被拒绝" 'echo "$BODY" | grep -q "账号密码不正确"' "$BODY" "$STATUS"

echo ""

# ==========================================
# 第二部分：用户管理
# ==========================================
echo -e "${YELLOW}[第二部分] 用户管理${NC}"
echo "---"
AUTH="Authorization: Bearer $ACCESS_TOKEN"

echo "  2.1 GET /user/detail"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/detail" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "获取用户详情" "获取成功" "$BODY" "$STATUS"
USER_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
if [ -z "$USER_ID" ]; then
    RESP=$(curl -s "$BASE_URL/user/listPage?current_page=1&page_size=10" -H "$AUTH")
    USER_ID=$(echo "$RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

echo "  2.2 GET /user/listPage"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/listPage?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "用户分页查询" "获取成功" "$BODY" "$STATUS"

echo "  2.3 POST /user/update"
if [ -n "$USER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/user/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$USER_ID\",\"name\":\"更新后的名字\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "用户更新" "操作成功" "$BODY" "$STATUS"
else
    report_skip "用户更新 (无ID)"
fi

echo "  2.4 POST /user/resetPassword"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/user/resetPassword" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"mobile":"13800138001","old_password":"Test12345678","new_password":"NewPass12345"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "重置密码" "密码重置成功" "$BODY" "$STATUS"

echo "  2.5 GET /user/getUserPermission"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/getUserPermission" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "获取用户权限" "获取成功" "$BODY" "$STATUS"

echo "  2.6 POST /user/del"
if [ -n "$USER_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/user/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":[\"$USER_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "用户删除" "删除成功" "$BODY" "$STATUS"
else
    report_skip "用户删除 (无ID)"
fi

echo ""

# ==========================================
# 第三部分：角色管理 (utils.Success(c) → "操作成功")
# ==========================================
echo -e "${YELLOW}[第三部分] 角色管理${NC}"
echo "---"

echo "  3.1 POST /role/add"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/role/add" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"role_name":"测试角色","role_code":"test_role","status":"01","description":"测试用角色"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "创建角色" "操作成功" "$BODY" "$STATUS"

echo "  3.2 GET /role/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/role/list?name=测试" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "角色列表" "获取成功" "$BODY" "$STATUS"
ROLE_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  3.3 GET /role/listPage"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/role/listPage?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "角色分页查询" "获取成功" "$BODY" "$STATUS"

echo "  3.4 GET /role/detail"
if [ -n "$ROLE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/role/detail?id=$ROLE_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "角色详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "角色详情 (无ID)"
fi

echo "  3.5 POST /role/update"
if [ -n "$ROLE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/role/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$ROLE_ID\",\"role_name\":\"更新角色\",\"role_code\":\"updated_role\",\"status\":\"01\",\"description\":\"更新描述\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新角色" "操作成功" "$BODY" "$STATUS"
else
    report_skip "更新角色 (无ID)"
fi

echo "  3.6 POST /role/del"
if [ -n "$ROLE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/role/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":[\"$ROLE_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除角色" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除角色 (无ID)"
fi

echo ""

# ==========================================
# 第四部分：菜单管理 (utils.Success(c) → "操作成功")
# ==========================================
echo -e "${YELLOW}[第四部分] 菜单管理${NC}"
echo "---"

echo "  4.1 POST /menu/add"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/menu/add" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"parent_id":"0","name":"测试菜单","type":1,"perms":"test:menu","status":1,"sort":99,"icon":"test","description":"测试菜单"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "创建菜单" "操作成功" "$BODY" "$STATUS"

echo "  4.2 GET /menu/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/menu/list" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "菜单列表" "获取成功" "$BODY" "$STATUS"
MENU_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  4.3 GET /menu/listPage"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/menu/listPage?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "菜单分页查询" "获取成功" "$BODY" "$STATUS"

echo "  4.4 GET /menu/detail"
if [ -n "$MENU_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/menu/detail?id=$MENU_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "菜单详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "菜单详情 (无ID)"
fi

echo "  4.5 POST /menu/update"
if [ -n "$MENU_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/menu/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$MENU_ID\",\"name\":\"更新菜单\",\"parent_id\":\"0\",\"type\":1,\"status\":1}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新菜单" "操作成功" "$BODY" "$STATUS"
else
    report_skip "更新菜单 (无ID)"
fi

echo "  4.6 POST /menu/del"
if [ -n "$MENU_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/menu/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"ids\":[\"$MENU_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除菜单" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除菜单 (无ID)"
fi

echo ""

# ==========================================
# 第五部分：操作日志
# ==========================================
echo -e "${YELLOW}[第五部分] 操作日志${NC}"
echo "---"

echo "  5.1 POST /OperationLogList"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/OperationLogList?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "操作日志分页查询" "获取成功" "$BODY" "$STATUS"

echo ""

# ==========================================
# 第六部分：报修管理
# ==========================================
echo -e "${YELLOW}[第六部分] 报修管理${NC}"
echo "---"

echo "  6.1 POST /repair/submit"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/repair/submit" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"repair_type":1,"address":"测试楼101","description":"测试报修问题","images":[],"contact":"张三","phone":"13800138001"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "提交报修" "提交成功" "$BODY" "$STATUS"

echo "  6.2 GET /repair/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/repair/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "报修列表" "获取成功" "$BODY" "$STATUS"
REPAIR_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  6.3 GET /repair/detail"
if [ -n "$REPAIR_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/repair/detail?id=$REPAIR_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "报修详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "报修详情 (无ID)"
fi

echo "  6.4 POST /repair/update"
if [ -n "$REPAIR_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/repair/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$REPAIR_ID\",\"repair_type\":1,\"address\":\"测试楼101\",\"contact\":\"张三\",\"phone\":\"13800138001\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新报修" "更新成功" "$BODY" "$STATUS"
else
    report_skip "更新报修 (无ID)"
fi

echo "  6.5 POST /repair/record"
if [ -n "$REPAIR_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/repair/record" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$REPAIR_ID\",\"status\":2,\"remark\":\"测试处理记录\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "报修处理记录" "状态更新成功" "$BODY" "$STATUS"
else
    report_skip "报修处理记录 (无ID)"
fi

# 6.6 报修删除使用 query ?id= 而非 JSON body
echo "  6.6 POST /repair/del"
if [ -n "$REPAIR_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/repair/del?id=$REPAIR_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除报修" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除报修 (无ID)"
fi

echo ""

# ==========================================
# 第七部分：校区管理 (返回 "创建成功"/"更新成功" 等自定义消息)
# ==========================================
echo -e "${YELLOW}[第七部分] 校区管理${NC}"
echo "---"

echo "  7.1 POST /campus/create"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/campus/create" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"campus_name":"测试校区","address":"测试地址","contact":"张三","phone":"010-12345678"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "创建校区" "创建成功" "$BODY" "$STATUS"

echo "  7.2 GET /campus/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/campus/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "校区列表" "获取成功" "$BODY" "$STATUS"
CAMPUS_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  7.3 GET /campus/detail"
if [ -n "$CAMPUS_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/campus/detail?id=$CAMPUS_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "校区详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "校区详情 (无ID)"
fi

echo "  7.4 GET /campus/all"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/campus/all" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "所有校区" "获取成功" "$BODY" "$STATUS"

echo "  7.5 POST /campus/update"
if [ -n "$CAMPUS_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/campus/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$CAMPUS_ID\",\"campus_name\":\"更新校区\",\"address\":\"新地址\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新校区" "更新成功" "$BODY" "$STATUS"
else
    report_skip "更新校区 (无ID)"
fi

# 校区删除使用 ids (不是 id)
echo "  7.6 POST /campus/del"
if [ -n "$CAMPUS_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/campus/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"ids\":[\"$CAMPUS_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除校区" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除校区 (无ID)"
fi

echo ""

# ==========================================
# 第八部分：楼栋管理
# ==========================================
echo -e "${YELLOW}[第八部分] 楼栋管理${NC}"
echo "---"

RESP=$(curl -s -X POST "$BASE_URL/campus/create" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"campus_name":"关联校区","address":"关联地址","contact":"联系人","phone":"010-88888888"}')
RESP2=$(curl -s "$BASE_URL/campus/all" -H "$AUTH")
CAMPUS_ID2=$(echo "$RESP2" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  8.1 POST /building/create"
if [ -n "$CAMPUS_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/building/create" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"campus_id\":\"$CAMPUS_ID2\",\"building_no\":\"B001\",\"building_name\":\"测试楼\",\"floor_count\":5,\"room_count\":20}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "创建楼栋" "创建成功" "$BODY" "$STATUS"
else
    report_skip "创建楼栋 (无校区ID)"
fi

echo "  8.2 GET /building/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/building/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "楼栋列表" "获取成功" "$BODY" "$STATUS"
BUILDING_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  8.3 GET /building/detail"
if [ -n "$BUILDING_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/building/detail?id=$BUILDING_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "楼栋详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "楼栋详情 (无ID)"
fi

echo "  8.4 GET /building/byCampus"
if [ -n "$CAMPUS_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/building/byCampus?campus_id=$CAMPUS_ID2" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "按校区查询楼栋" "获取成功" "$BODY" "$STATUS"
else
    report_skip "按校区查询楼栋 (无校区ID)"
fi

echo "  8.5 POST /building/update"
if [ -n "$BUILDING_ID" ] && [ -n "$CAMPUS_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/building/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$BUILDING_ID\",\"campus_id\":\"$CAMPUS_ID2\",\"building_no\":\"B001\",\"building_name\":\"更新楼\",\"floor_count\":6,\"room_count\":25}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新楼栋" "更新成功" "$BODY" "$STATUS"
else
    report_skip "更新楼栋 (缺少ID)"
fi

# 楼栋删除使用 ids
echo "  8.6 POST /building/del"
if [ -n "$BUILDING_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/building/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"ids\":[\"$BUILDING_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除楼栋" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除楼栋 (无ID)"
fi

echo ""

# ==========================================
# 第九部分：宿舍管理
# ==========================================
echo -e "${YELLOW}[第九部分] 宿舍管理${NC}"
echo "---"

RESP=$(curl -s -X POST "$BASE_URL/building/create" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d "{\"campus_id\":\"$CAMPUS_ID2\",\"building_no\":\"B002\",\"building_name\":\"宿舍楼\",\"floor_count\":6,\"room_count\":30}")
RESP=$(curl -s "$BASE_URL/building/list?current_page=1&page_size=10" -H "$AUTH")
BUILDING_ID2=$(echo "$RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  9.1 POST /dorm/create"
if [ -n "$BUILDING_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/dorm/create" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"building_id\":\"$BUILDING_ID2\",\"room_no\":\"101\",\"floor\":1,\"room_type\":1,\"max_count\":4}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "创建宿舍" "创建成功" "$BODY" "$STATUS"
else
    report_skip "创建宿舍 (无楼栋ID)"
fi

echo "  9.2 GET /dorm/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/dorm/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "宿舍列表" "获取成功" "$BODY" "$STATUS"
DORM_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  9.3 GET /dorm/detail"
if [ -n "$DORM_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/dorm/detail?id=$DORM_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "宿舍详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "宿舍详情 (无ID)"
fi

echo "  9.4 POST /dorm/update"
if [ -n "$DORM_ID" ] && [ -n "$BUILDING_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/dorm/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$DORM_ID\",\"building_id\":\"$BUILDING_ID2\",\"room_no\":\"102\",\"floor\":1,\"room_type\":2,\"max_count\":6}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新宿舍" "更新成功" "$BODY" "$STATUS"
else
    report_skip "更新宿舍 (缺少ID)"
fi

echo "  9.5 GET /dorm/warning"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/dorm/warning" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "宿舍容量预警" "获取成功" "$BODY" "$STATUS"

# 宿舍删除使用 ids
echo "  9.6 POST /dorm/del"
if [ -n "$DORM_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/dorm/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"ids\":[\"$DORM_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除宿舍" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除宿舍 (无ID)"
fi

echo ""

# ==========================================
# 第十部分：水电费管理
# ==========================================
echo -e "${YELLOW}[第十部分] 水电费管理${NC}"
echo "---"

RESP=$(curl -s "$BASE_URL/dorm/list?current_page=1&page_size=10" -H "$AUTH")
DORM_ID2=$(echo "$RESP" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  10.1 GET /utility/price"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/utility/price" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "获取价格配置" "获取成功" "$BODY" "$STATUS"

echo "  10.2 POST /utility/create"
if [ -n "$DORM_ID2" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/utility/create" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"room_id\":\"$DORM_ID2\",\"year\":2026,\"month\":5,\"water_usage\":5.5,\"electric_usage\":30.2}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "创建水电费记录" "创建成功" "$BODY" "$STATUS"
else
    report_skip "创建水电费记录 (无宿舍ID)"
fi

echo "  10.3 GET /utility/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/utility/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "水电费列表" "获取成功" "$BODY" "$STATUS"
UTILITY_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  10.4 GET /utility/detail"
if [ -n "$UTILITY_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/utility/detail?id=$UTILITY_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "水电费详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "水电费详情 (无ID)"
fi

echo "  10.5 POST /utility/pay"
if [ -n "$UTILITY_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/utility/pay" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$UTILITY_ID\"}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "水电费缴费" "缴费成功" "$BODY" "$STATUS"
else
    report_skip "水电费缴费 (无ID)"
fi

echo "  10.6 POST /utility/price"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/utility/price" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"water_price":3.5,"electric_price":0.6}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "更新价格" "更新成功" "$BODY" "$STATUS"

echo "  10.7 GET /utility/statistics"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/utility/statistics" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "水电费统计" "获取成功" "$BODY" "$STATUS"

echo "  10.8 GET /utility/warning"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/utility/warning" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "欠费预警" "获取成功" "$BODY" "$STATUS"

echo ""

# ==========================================
# 第十一部分：公告管理
# ==========================================
echo -e "${YELLOW}[第十一部分] 公告管理${NC}"
echo "---"

echo "  11.1 POST /notice/create"
RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/notice/create" \
    -H "$AUTH" -H "Content-Type: application/json" \
    -d '{"title":"测试公告","content":"这是一个测试公告内容","notice_type":1,"is_top":2,"publish_time":"2026-06-06T10:00:00Z"}')
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "创建公告" "创建成功" "$BODY" "$STATUS"

echo "  11.2 GET /notice/list"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/notice/list?current_page=1&page_size=10" -H "$AUTH")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
report "公告列表" "获取成功" "$BODY" "$STATUS"
NOTICE_ID=$(echo "$BODY" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

echo "  11.3 GET /notice/detail"
if [ -n "$NOTICE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/notice/detail?id=$NOTICE_ID" -H "$AUTH")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "公告详情" "获取成功" "$BODY" "$STATUS"
else
    report_skip "公告详情 (无ID)"
fi

echo "  11.4 POST /notice/update"
if [ -n "$NOTICE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/notice/update" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$NOTICE_ID\",\"title\":\"更新公告\",\"content\":\"更新内容\",\"notice_type\":1,\"is_top\":2}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "更新公告" "更新成功" "$BODY" "$STATUS"
else
    report_skip "更新公告 (无ID)"
fi

echo "  11.5 POST /notice/top"
if [ -n "$NOTICE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/notice/top" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"id\":\"$NOTICE_ID\",\"is_top\":1}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "置顶公告" "置顶设置成功" "$BODY" "$STATUS"
else
    report_skip "置顶公告 (无ID)"
fi

# 公告删除使用 ids
echo "  11.6 POST /notice/del"
if [ -n "$NOTICE_ID" ]; then
    RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/notice/del" \
        -H "$AUTH" -H "Content-Type: application/json" \
        -d "{\"ids\":[\"$NOTICE_ID\"]}")
    STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
    report "删除公告" "删除成功" "$BODY" "$STATUS"
else
    report_skip "删除公告 (无ID)"
fi

echo ""

# ==========================================
# 第十二部分：鉴权测试
# ==========================================
echo -e "${YELLOW}[第十二部分] 鉴权测试${NC}"
echo "---"

echo "  12.1 无 Token 访问受保护接口"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/listPage")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "无 Token 被正确拒绝" '[ "$STATUS" = "401" ]' "$BODY" "$STATUS"

echo "  12.2 无效 Token"
RESP=$(curl -s -w "\n%{http_code}" "$BASE_URL/user/listPage" -H "Authorization: Bearer invalid.token.here")
STATUS=$(echo "$RESP" | tail -1); BODY=$(echo "$RESP" | sed '$d')
check "无效 Token 被正确拒绝" '[ "$STATUS" = "401" ]' "$BODY" "$STATUS"

echo ""

# ==========================================
# 汇总
# ==========================================
echo "========================================"
echo "           测试结果汇总"
echo "========================================"
echo ""
echo "  总用例数:   $TOTAL"
echo -e "  通过:       ${GREEN}$PASS${NC}"
echo -e "  失败:       ${RED}$FAIL${NC}"
echo -e "  跳过:       ${BLUE}$SKIP${NC}"
echo ""
RATE=0
if [ "$TOTAL" -gt 0 ]; then
    RATE=$(awk "BEGIN {printf \"%.1f\", ($PASS/$TOTAL)*100}")
fi
echo "  通过率:      $RATE%"
echo ""
if [ "$FAIL" -eq 0 ]; then
    echo -e "  ${GREEN}✓ 所有测试通过！${NC}"
else
    echo -e "  ${RED}✗ 存在 $FAIL 个失败用例${NC}"
fi
echo ""
