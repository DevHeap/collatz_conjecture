## Collatz Conjecture playground webapp

[![Build Status](https://travis-ci.org/DevHeap/collatz_conjecture.svg?branch=dev)](https://travis-ci.org/DevHeap/collatz_conjecture)

### What is it all about?
Let us take any positive integer, then check if it is even or odd.
* If it is even, divide it by 2
* If it is odd, multiply it by 3 and then add 1

Now take the number you got and repeat the steps. Over and over and over. Eventually, you will get 1, according to the __*Collatz conjecture*__. Or will you? Nobody has proved it yet.

However, the conjecture has been checked for all numbers up to 87×2^60.

_Example:_ Let us take 42 as a starting number. Then we get 42->21->64->32->16->8->4->2->1 as a __*Collatz path*__.

**OK, got it. So what is your software doing?**  
Well, it is all pretty simple. You enter a number and click **Start**. We count a Collatz path starting with this number for you. And then for the next one. And then again for the next, and so on, infinitely or until you click **Stop**.  
In order for you not to get _too_ bored, we also show you some information:
* For every number:
  * Path length
  * Maximal number in the path
  * Average of numbers in the path
  * Time spent calculating the path
* Overall -- a histogram showing how many paths have every specified length.

All numbers are neatly sorted for your convenience. Also, you can enter arbitrarily long numbers there, as long as they are positive integers.

### How to deploy it?

We are using `docker-compose` for deployment. While being in the repo root folder, type:
```
$ docker-compose up -d
```

### How to contact developers?

If you _ever_ happen to have any questions, write us using emails specified in our GitHub accounts or open an issue.