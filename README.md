# mini-oauth2-proxy

[oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy)

mini-oauth2-proxyは、OAuth2認証を使用してアプリケーションやサービスへのアクセスを保護するためのリバースプロキシです。oauth2-proxyから派生したプロジェクトですが、よりシンプルで読みやすいコードを提供することを目的としています。

# なんのためのソフトウェアか
oauth2-proxyは、OAuth2認証を使用してアプリケーションやサービスへのアクセスを保護するための優れたソフトウェアです。しかし、時間とともに複雑化が進み、いくつかの問題点が明らかになってきました。

- 複雑な設定ファイル：oauth2-proxyの設定ファイルは非常に複雑で、理解するのが難しくなっています。多くのオプションがありますが、ドキュメントが不十分であり、それらがどこで使われているのか、どのような意味があるのかを把握するのは困難です。
- コードの肥大化：oauth2-proxyのコードは1万行以上に及び、メンテナンスが困難になっています。コードの理解に多くの時間を要するようになりました。
- OAuth2とは直接関係のない機能：Basic認証によるログインという、OAuth2とは全く関係のない機能が実装されています。これらの機能はoauth2-proxyの本来の目的から外れており、コードの複雑化を招いています。
- 設計上の問題点：いくつかの設計上の失敗が見られます。不要なインターフェースの存在や、責務の分離が不十分なところがあります。これらは、コードの理解を難しくし、保守性を低下させています。
- TLSによる暗号化処理：oauth2-proxyはTLSによる暗号化処理を内包しています。しかし、TLSはOAuth2とは直接的に関係がなく、他の要素（nginxなど）に任せるべき責務です。
- 無理やりにステートレスであること：oauth2-proxyはID TokenをCookieに保存するかRedisと連携して保存することによって、状態の保持を避けています。しかし、Cookieへの保存は不要なIDTokenの暗号化や署名などの問題およびCookieの大きさの制限による不要な複雑化を招いています。また、Redisへの保存は簡易的なユースケースでは不要であり、サーバー管理の複雑化を引き起こします。

これらの問題点を解消し、OIDC認証のためのシンプルで読みやすいリバースプロキシを提供したいという思いから、mini-oauth2-proxyの開発が始まりました。
mini-oauth2-proxyは、oauth2-proxyの本質的な機能に集中し、oauth2-proxyよりも大幅に少ないコード量で実装されています。
また、注意深い設計によって、oauth2-proxyよりも読みやすく、メンテナンスしやすいコードになっています。

mini-oauth2-proxyがOAuth2認証のためのリバースプロキシを必要とする多くの開発者にとって、優れた選択肢になると信じています。
さらに、mini-oauth2-proxyは非常にシンプルなOIDCのRPの実装を提供しているため、教育用にも適しています。

# 注意点

- テストが無い：まだテストがありません。追加したいと思っています。
- 運用の実績がない：mini-oauth2-proxyは、oauth2-proxyの動作を理解するための学習用とでの使用と、開発段階における簡易的なoauth2-proxyの代替としての使用を想定しています。 **本番環境で使わないでください。**

# 設定

簡易的な設定ファイルを示します。ファイル名は`mini-oauth2-proxy.json`である必要があります。
```json
{
    "oidc" : {
        "providers": [
            {
                "id": "IdP ID",
                "clientID": "CLIENT ID",
                "clientSecret": "CLIENT SECRET",
                "redirectURL": "REDIRECT URL",
                "startPath": "/start",
                "scopes": [
                    "SCOPE1", "SCOPE2"
                ],
                "issuer": "ISSUER URL"
            }
        ],
        "skipLoginPage": true
    },
    "upstream" : {
        "servers": [
            {
                "id": "myservice",
                "url": "http://localhost:3000",
                "matchPath": "/myservice"
            }
        ]
    },
    "headerInjection": {
        "request": [
            {
                "name": "X-Authenticated-User",
                "type": "idTokenClaim",
                "values": [
                    "name"
                ]
            },
            {
                "name": "X-Authenticated-EMail",
                "type": "userInfo",
                "values": [
                    "email"
                ]
            }
        ],
        "response": []
    },
    "proxyURL": {
        "host": "mini-oauth2-proxyよりもインターネット側で動作するTLS終端を行うプロキシのホスト名"
    },
    "log": {
        "level": "Info"
    },
    "port": 8080
}
```

# Contribution

プルリクエストや Issue は大歓迎です。mini-oauth2-proxy をより良いものにするために、ぜひご協力ください。
