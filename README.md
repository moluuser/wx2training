# wx2training
将使用[WechatExporter](https://github.com/BlueMatthew/WechatExporter)导出的微信聊天记录转换为可供模型训练或微调的数据集格式。
* 在`main`分支下可导出带有`history`的供[ChatGLM-6B ptuning](https://github.com/THUDM/ChatGLM-6B/tree/main/ptuning)微调模型的对话数据。
* 在`lora`分支下可导出供对ChatGLM模型使用lora微调方案的数据集格式(`{"instruction":"","input":"","output":""}`)。
* 需要对WechatExporter导出的文本数据进行预处理，比如将`「`和`」`删除。
