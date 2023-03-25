# Azure OpenAI Proxy

## 介绍

<a href="./README.md">English</a> |
<a href="./README.zh-cn.md">中文</a>

一个Azure OpenAI API的代理工具，可以将一个OpenAI请求转化为Azure
OpenAI请求，方便作为各类开源ChatGPT的后端使用。同时也支持作为单纯的OpenAI接口代理使用，用来解决OpenAI接口在部分地区的被限制使用的问题。

亮点：

+ 🌐 支持代理所有 Azure OpenAI 接口
+ 🧠 支持代理所有 Azure OpenAI 模型以及自定义微调模型
+ 🗺️ 支持自定义 Azure 部署名与 OpenAI 模型的映射关系
+ 🔄 支持反向代理和正向代理两种方式使用

## 使用方式

### 1. 作为反向代理使用（即一个OpenAI API网关）

环境变量

| 参数名                        | 描述                                                                                                                                                                               | 默认值                                                                     |
|:---------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:------------------------------------------------------------------------|
| AZURE_OPENAI_PROXY_ADDRESS | 服务监听地址                                                                                                                                                                           | 0.0.0.0:8080                                                            |
| AZURE_OPENAI_PROXY_MODE    | 代理模式，可以为azure/openai 2种模式                                                                                                                                                        | azure                                                                   |
| AZURE_OPENAI_ENDPOINT      | Azure OpenAI Endpoint，一般类似https://{custom}.openai.azure.com的格式。必需。                                                                                                               |                                                                         |
| AZURE_OPENAI_APIVERSION    | Azure OpenAI API 的 API 版本。默认为 2023-03-15-preview。                                                                                                                                | 2023-03-15-preview                                                      |
| AZURE_OPENAI_MODEL_MAPPER  | 一个逗号分隔的 model=deployment 对列表。模型名称映射到部署名称。例如，`gpt-3.5-turbo=gpt-35-turbo`,`gpt-3.5-turbo-0301=gpt-35-turbo-0301`。未匹配到的情况下，代理会直接透传model作为deployment使用(其实Azure大部分模型名字和OpenAI的保持一致)。 | `gpt-3.5-turbo=gpt-35-turbo`<br/>`gpt-3.5-turbo-0301=gpt-35-turbo-0301` |

在命令行调用

```shell
curl https://{your-custom-domain}/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {your azure api key}" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

```

### 2. 作为正向代理使用（即一个HTTP Proxy）

通过HTTP访问Azure OpenAI接口时，可以直接作为代理使用，但是这个工具没有内置原生的HTTPS支持，需要在工具前架设一个类似Nginx的HTTPS代理，来支持访问HTTPS版本的OpenAI接口。

假设你配置好后的代理域名为`https://{your-domain}.com`，你可以在终端中执行以下命令来配置http代理：

```shell
export https_proxy=https://{your-domain}.com 

curl https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {your azure api key}" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'

```

或者在其他开源Web ChatGPT项目中配置为HTTP代理

```
export HTTPS_PROXY=https://{your-domain}.com
```

## 部署方式

通过docker部署

```shell
docker pull ishadows/azure-openai-proxy:latest
docker run -d -p 8080:8080 --name=azure-openai-proxy \
  --env AZURE_OPENAI_ENDPOINT={your azure endpoint} \
  --env AZURE_OPENAI_MODEL_MAPPER={your custom model mapper ,like: gpt-3.5-turbo=gpt-35-turbo,gpt-3.5-turbo-0301=gpt-35-turbo-0301} \
  ishadows/azure-openai-proxy:latest
```
调用
```shell
curl https://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {your azure api key}" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

## 模型映射机制

`AZURE_OPENAI_MODEL_MAPPER`中预定义了一系列模型映射的规则，默认配置基本上满足了所有Azure模型的映射，规则包括：

+ `gpt-3.5-turbo` -> `gpt-35-turbo`
+ `gpt-3.5-turbo-0301` -> `gpt-35-turbo-0301`
+ 以及一个透传模型名的机制作为fallback手段

对于自定义的微调模型，可以直接透传模型名。对于部署名字和模型名不一样的，可以自定义映射关系，比如：

| 模型名称               | 部署名称                         |
|:-------------------|:-----------------------------|
| gpt-3.5-turbo      | gpt-35-turbo-upgrade         |
| gpt-3.5-turbo-0301 | gpt-35-turbo-0301-fine-tuned |

## 许可证

MIT