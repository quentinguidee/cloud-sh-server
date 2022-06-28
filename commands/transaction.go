package commands

type Transaction struct {
	commands []Command
}

func NewTransaction(commands []Command) Transaction {
	return Transaction{commands}
}

func (t Transaction) Run() ICommandError {
	for i, command := range t.commands {
		err := command.Run()
		if err != nil {
			for j := i - 1; j >= 0; j-- {
				revertError := t.commands[j].Revert()
				if revertError != nil {
					return err
				}
			}
			return err
		}
	}
	return nil
}
