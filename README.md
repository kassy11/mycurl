### curlコマンドの実装

次の機能を満たすアプリケーションを作る

(1) curl https://example.com 相当のことができる機能  
(2) curl -o file https://example.com 相当のことができる機能  
(3) curl -v  https://example.com 相当のことができる機能  
(4) curl -X POST https://example.com 相当のことができる機能  
(5) curl -X POST -d "key=value" https://example.com 相当のことができる機能  

コマンドオプションは (1) ~ (5) に記載された組み合わせだけではなく、
任意の組み合わせができるようにする。