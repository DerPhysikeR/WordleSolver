# WordleSolver

This is a solver for the popular word game [Wordle](https://www.powerlanguage.co.uk/wordle/).

## How to use the solver

Download the executable for your platform and run it by double clicking it.
This is a command line application, so a CLI should automatically open.

You should be met with this screen:

```
Calculate best guesses ...
```

Now every word in the word list is scored against each other word, to find the best
possible guess.
This can take a few seconds (~10s on my machine).

```
Calculate best guesses from given word list ...
Best guesses: [AESIR REAIS SERAI AIERY AYRIE ARIEL RAILE ALOES REALO STOAE ANOLE AEROS]
Your guess:
```

Now choose one of the guesses, and type it into wordle.
Wordle then gives you a score for your guess (green = [H]it, yellow = [h]it, grey = [.]).
Now input the guess and the score of your guess into the solver.

```
Calculate best guesses from given word list ...
Best guesses: [AESIR REAIS SERAI AIERY AYRIE ARIEL RAILE ALOES REALO STOAE ANOLE AEROS]
Your guess: AESIR
Score of the guess: .h.h.
Best guesses: [CLINE LINTY TINED ALINE ANILE CLINT CLIPE INCLE LIGNE LINCH LINED LINTS]
Your guess: CLINE
```

It immediately shows you the best subsequent guesses, to reduce the search space even
more.

```
Calculate best guesses from given word list ...
Best guesses: [AESIR REAIS SERAI AIERY AYRIE ARIEL RAILE ALOES REALO STOAE ANOLE AEROS]
Your guess: AESIR
Score of the guess: .h.h.
Best guesses: [CLINE LINTY TINED ALINE ANILE CLINT CLIPE INCLE LIGNE LINCH LINED LINTS]
Your guess: CLINE
Score of the guess: h.hhH
Best guesses: [CHEMO CHEMS CHEWS CHEWY HAEMS HAWMS HIEMS MAHWA MANEH MENSH MINCY MYNAH]
Your guess: CHEMO
Score of the guess: h.h..
The solution is: WINCE
```

Either one of those guesses is already the solution or the solver shows you the last
possible remaining word at the end.


## How it works

It is basically a simplified version of
[Minimax](https://en.wikipedia.org/wiki/Minimax).

For every word (guess) a hashtable is created by iterating over every other word
(solution) and scoring them against each other, e.g.,
`score(guess="AESIR", solution="WINCE") -> ".h.h."`.

The hashtable uses the score as the key and the number of times the score comes up as
the value.

The maximum value in this hashtable is then used as the weight of the guess.

To produce the list of best guesses, the words are then sorted by their weights.

This effectively sorts the words, by how much they reduce the search space in the worst
case.

After each guess, the possible solution words are filtered until only one word is left,
which has to be the solution.

Unfortunately, this approach has a time complexity of `O(n^2)`, which is why the search
for the solution takes several seconds.
