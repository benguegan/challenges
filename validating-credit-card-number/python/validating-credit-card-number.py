# Enter your code here. Read input from STDIN. Print output to STDOUT
import sys
import re


def main():
    start_with_456 = r"(?=[456])"
    grouped_by_4 = r"(?=(\d{4}-?){3}\d{4})"
    no_more_than_4_consecutive_digits = r"(?!.*(\d)(-?\1){3})"
    contains_16_digits = r"(?=^(\d-?){16}$)"

    credit_card_pattern = re.compile(
        f"^{start_with_456}{no_more_than_4_consecutive_digits}{grouped_by_4}{contains_16_digits}.*?$"
    )

    for idx, line in enumerate(sys.stdin):
        if idx == 0:
            continue

        credit_card = line.strip()

        if credit_card_pattern.search(credit_card):
            print("Valid")
        else:
            print("Invalid")


if __name__ == "__main__":
    main()
