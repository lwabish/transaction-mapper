# transaction-mapper

将从各银行导出的交易明细转换成记账软件的模板，方便导入记账。

## 功能

1. 将各家银行导出的交易明细转换为指定记账APP的导入模板
2. 自动按照规则设置账户和分类

目前支持的银行:

- [中国工商银行](https://mybank.icbc.com.cn/icbc/newperbank/perbank3/frame/frame_index.jsp)
- [招商银行信用卡](https://www.cmbchina.com/)
- [招商银行储蓄卡](https://www.cmbchina.com/)（招行储蓄卡和信用卡导出的途径和格式不一样，需要分开解析）

目前支持的记账app:

- [Bluecoins](https://www.bluecoinsapp.com/)
- [钱迹](https://www.qianjiapp.com/)

## 使用方式

### 命令行模式
1. 从银行导出交易明细文件
2. 运行命令，完成转换(`go run main.go -h`查看完整命令帮助)
    ```bash
    # 招行信用卡(cmbCredit+bluecoins)
    go run main.go \
      -i cmbCredit-2501.json \
      -b cmbCredit \
      -a bluecoins \
      -z 信用卡 \
      -t 【信用卡】招商银行
    # 更多组合：
    # 工行信用卡/储蓄卡(icbc+bluecoins)
    # 招行储蓄卡(cmb+bluecoins)
    ```
3. 将生成的模板导入到记账app中

### web模式

todo

## 关于记账

1. 记账的难点是繁琐不易坚持，手工记账更是对体力和脑力的极大消耗
2. 早些年网易有钱的出现为记账自动化提供了很棒的路线和实现方式，可惜的是该产品未能持续维护。另外国内银行的数据开放程度比较低，为自动化带来了重重阻碍。
3. 我个人在记账自动化实践过一个比较长期的方案是用bluecoins的通知读取功能，读取工商银行app的动账通知，实现半自动记账
4. 但是上述方案也有一些弊端：比如依赖安卓平台、依赖通知的稳定性、每条交易需要随通知及时处理等
5. 因此考虑改用目前各记账APP支持更好的模板导入模式，于是有了本工具进行辅助

# todo

- ✅支持解析转账类交易
- ✅支持钱迹
- ✅cobra flag parser & multi cmd
- ✅账户信息页从flag中解析
- ✅config解析优化
- ✅server模式
- ✅前端界面
- ✅ci
- 文档：各银行导出的方法
- 支持跳过交易
- release image/基于容器的使用说明
- 多用户(接入logto)
- 基于用户的，可维护配置的config/账户信息

# prompt

使用LLM快速针对不同的结构生成结构体，prompt：

```text
我将提供你一个用逗号分隔的csv标题行，请将每一个csv列对应go结构体的一个字段，示例标题行及结构体如下
示例标题行：交易日期,记账日期,摘要,交易场所,交易国家或地区简称,交易金额(收入),交易金额(支出)
示例结构体：
type icbcTxn struct {
	TranDate             string `csv:"交易日期"`
	AccountDate          string `csv:"记账日期"`
	Abstract             string `csv:"摘要"`
	Platform             string `csv:"交易场所"`
	CountryRegion        string `csv:"交易国家或地区简称"`
	TranAmountIncome     string `csv:"交易金额(收入)"`
	TranAmountOutcome    string `csv:"交易金额(支出)"`
	TranCurrency         string `csv:"交易币种"`
	AccountAmountIncome  string `csv:"记账金额(收入)"`
	AccountAmountOutcome string `csv:"记账金额(支出)"`
	AccountCurrency      string `csv:"记账币种"`
	Balance              string `csv:"余额"`
	CounterpartyName     string `csv:"对方户名"`
注意：结构体的字段名称请根据标题行翻译成准确的英文，使用首字母大写的驼峰形式，类型全部是string
```

```text
我将给你提供一个json文件，请根据该json结构体，生成用于反序列化的go结构体。字段名称全部转换为首字母大写的驼峰形式，类型全部是string。
{
    "credit_limit": "¥ 62,000.00",
    "payment_due_date": "2025年05月08日",
    "current_balance": "¥ 4,409.63",
    "minimum_payment": "¥ 220.48",
    "statement_date": "2025年04月20日",
    "transaction_details": [
        {
            "sold_date": "04/08",
            "posted_date": "04/08",
            "description": "自动还款",
            "rmb_amount": "-3,859.80",
            "card_no": "7139",
            "original_tran_amount": "-3,859.80"
        }
    ],
    "current_balance_summary": "¥ 4,409.63",
    "balance_b_f": "¥ 3,859.80",
    "payment": "¥ 3,859.80",
    "new_charges": "¥ 4,738.97",
    "adjustment": "¥ 329.34",
    "interest": "¥ 0.00"
}
```