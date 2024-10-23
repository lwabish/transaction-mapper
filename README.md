# transaction-mapper

将从各银行导出的交易明细转换成记账软件的模板，以供导入记账

目前支持的银行:
- [中国工商银行](https://mybank.icbc.com.cn/icbc/newperbank/perbank3/frame/frame_index.jsp)
- [招商银行](https://www.cmbchina.com/)

目前支持的记账app:
- [Bluecoins](https://www.bluecoinsapp.com/)

# usage

1. 从银行导出交易明细文件
2. 按照如下规则重命名，主要为了设定导入账户： `$bankType-$account-$subaccount-other.csv`
3. 运行命令，完成转换
    ```commandline
    # 工行信用卡/储蓄卡(icbc+bluecoins)
    go run cmd/main.go -consumer bluecoins -input icbc-信用卡-【信用卡】工商银行-240714.csv 
    
    # 招行信用卡(cmbCredit+bluecoins)
    go run ./cmd/main.go -consumer bluecoins -input ./cmbCredit-信用卡-【信用卡】招商银行-2407.json
    
    # 招行储蓄卡(cmb+bluecoins)
    go run ./cmd/main.go -consumer bluecoins -input ./cmb-现金-【借记卡】招商银行-2407.csv
    ```
4. 将生成的模板导入到记账app中

# todo

1. 识别转账类交易
2. 更优雅的账户设置方案
3. 支持分类自定义
4. 增强自动分类匹配
5. 支持钱迹
6. 文档：各银行导出
7. update flag parser
8. 前端界面
