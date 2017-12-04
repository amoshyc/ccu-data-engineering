# Final Project Proposal

403410034 資工四 黃鈺程

我想了兩個題目，之後會選一個來實作，希望題目沒有太過離題。

## 題目 1: 藏頭詩

我半年前看到的[投影片](https://www.slideshare.net/ckmarkohchang/ss-45210079)，實作了這個題目，我覺得非常有趣，想來實作看看。例如該投影片中的一個結果為：

![](https://i.imgur.com/V4c4fhG.png)

該投影片中提到了許多方法，包含 Viterbi(Hidden Markov Model) 與 RNN，並討論到了時間複雜度，想讓演算法快且正確是一個很難的問題。

## 題目 2: 利用 GAN 生成 Haiku

題目發想來自於昨天（12/2），霍金為了慶祝 FB 追蹤人數來到 4,000,000 人，[po 了一篇文](https://www.facebook.com/stephenhawking/posts/1600494156704342)請大家以 Science 為題創作 haiku（徘句），並從中選出一人贈送他的著作 <時間簡史>。英文的徘句即為三個句子，音節數分別為 5-7-5，例如:

```
Stars, light-years away
Vanished, but their brilliance stays
Life is, memories.
```

我想試試能不能使用 Attention + GAN 之類的方法生成 haiku，預期這會是一個非常難的問題，畢竟我沒有 NLP 的基礎。