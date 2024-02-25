package redirect

// Optionから、指定されたリクエストに対応するupstreamのリダイレクトパスを決定する
// 抽象化：どのようなアルゴリズムでリクエストからupstreamのパスを抽出するか
// 文脈：リクエストをoauth2proxyが補足したときに、upstreamの適切なリダイレクト先を設定する必要があるため、その処理を記述する
