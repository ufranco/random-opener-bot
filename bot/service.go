package bot

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Service interface {
	GetRandomOpener(message *discordgo.MessageCreate) (Opener, error)
	GetOpenerLeaderboard() ([]Opener, error)
	SetFavoriteOpener(message *discordgo.MessageCreate) (string, error)
	ProcessReaction(message *discordgo.MessageReactionAdd, reactedTo *discordgo.Message) error
}

type service struct {
	accountRepository AccountRepository
	openerRepository  OpenerRepository
}

func NewService(aRepository AccountRepository, oRepository OpenerRepository) Service {
	return &service{aRepository, oRepository}
}

func (logic *service) GetRandomOpener(message *discordgo.MessageCreate) (Opener, error) {

	opener, err := logic.openerRepository.GetRandomOpener()

	if err != nil {
		log.Fatal(err.Error())
		return Opener{}, err
	}

	return opener, nil
}

func (logic *service) GetOpenerLeaderboard() ([]Opener, error) {
	return logic.openerRepository.GetLeaderboard()
}

func (logic *service) SetFavoriteOpener(message *discordgo.MessageCreate) (string, error) {
	regexCompiler := regexp.MustCompile(`\!(randomOpener favorite|ro fav|rof)\s`)
	openerName := regexCompiler.ReplaceAllString(message.Content, "")
	openerName = strings.ToUpper(openerName)

	log.Printf("Opener name: %s", openerName)

	return openerName, logic.setFavoriteOpener(message.Author.ID, openerName)
}

func (logic *service) ProcessReaction(message *discordgo.MessageReactionAdd, reactedTo *discordgo.Message) error {
	return nil
	//return SetFavoriteOpener()
}

func (logic *service) setFavoriteOpener(accountId string, opener string) error {
	_, err := logic.openerRepository.FindById(opener)

	if err != nil && err.Error() == "mongo: no documents in result" {
		return fmt.Errorf("%s opener not found", opener)
	}

	account, err := logic.accountRepository.FindById(accountId)

	if err != nil && err.Error() == "mongo: no documents in result" {

		account = Account{
			ID:             accountId,
			FavoriteOpener: opener,
		}

		if err = logic.accountRepository.Register(account); err != nil {

			log.Fatal(err.Error())
			return err
		}
		if err = logic.openerRepository.UpdateReactionBy(opener, 1); err != nil {

			log.Fatal(err.Error())
			return err
		}

		return err

	} else if account.FavoriteOpener == opener {
		return fmt.Errorf("%s is already your favorite opener", opener)
	}

	if err = logic.openerRepository.UpdateReactionBy(account.FavoriteOpener, -1); err != nil {
		log.Fatal(err.Error())
		return err
	}

	if err = logic.openerRepository.UpdateReactionBy(opener, 1); err != nil {
		log.Fatal(err.Error())
		return err
	}

	err = logic.accountRepository.UpdateFavoriteOpener(accountId, opener)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}
