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

1. 从银行导出交易明细文件
2. 运行命令，完成转换(`go run main.go -h`查看完整命令帮助)
    ```bash
    # 招行信用卡(cmbCredit+bluecoins)
    go run main.go \
      -i cmbCredit-2501.json \
      -b cmbCredit \
      -c bluecoins \
      -z 信用卡 \
      -a 【信用卡】招商银行
    # 更多组合：
    # 工行信用卡/储蓄卡(icbc+bluecoins)
    # 招行储蓄卡(cmb+bluecoins)
    ```
3. 将生成的模板导入到记账app中

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
- server模式
- 前端界面
- 支持跳过某些对冲类交易
- 文档：各银行导出的方法
- release binary
- ci
- e2e test
