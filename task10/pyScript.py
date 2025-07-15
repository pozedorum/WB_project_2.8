#!/usr/bin/env python3

output_file = "tests/large_file.txt"
line_count = 100000
words_per_line = 10
word = ["word1","word8","word4","word2","word0","word9","word5","word3"]

print(f"Generating file with {line_count} lines...")
wordL = len(word)
with open(output_file, 'w') as f:
    for i in range(1, line_count + 1):
        line = ' '.join([word[i%wordL]] * words_per_line)
        f.write(line + '\n')
        
        if i % 1000 == 0:
            print(f"Generated {i} lines", end='\r')

print(f"\nDone. File created: {output_file}")