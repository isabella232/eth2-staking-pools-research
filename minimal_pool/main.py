import participant

def main():
    p = participant.Participant()
    p.generate_polynomial(1111,20)
    print(p.distribuite_shares([1,2,3,4,5,6,7,8,9]))

if __name__ == '__main__':
    main()
