		*=$C000
		LDA #0
loop	STX $D020
		DEX
		JMP loop
		JMP *
		BNE *

