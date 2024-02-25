package upstream

// 上流へのリクエストのdemuxを行う
// 抽象化：upstreamのpathのバリデーション、いい感じのrouterの構築
// 文脈依存度：oauth2proxyの初期化処理において、認証が成功したときにリクエストを送信するときに使うrouterを作成するためのモジュール
