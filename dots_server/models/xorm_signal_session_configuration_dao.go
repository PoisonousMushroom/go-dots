package models

import (
	"github.com/nttdots/go-dots/dots_server/db_models"
	log "github.com/sirupsen/logrus"
)

/*
 * Stores SignalSessionConfiguration to the DB.
 *
 * parameter:
 *  signalSessionConfiguration SignalSessionConfiguration
 *  customer Customer
 * return:
 *  err error
 */
func CreateSignalSessionConfiguration(signalSessionConfiguration SignalSessionConfiguration, customer Customer) (newSignalSessionConfiguration db_models.SignalSessionConfiguration, err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// same session_id data check
	dbSignalSessionConfiguration := new(db_models.SignalSessionConfiguration)
	_, err = engine.Where("customer_id = ? AND session_id = ?", customer.Id, signalSessionConfiguration.SessionId).Get(dbSignalSessionConfiguration)
	if err != nil {
		return
	}
	if dbSignalSessionConfiguration.Id != 0 {
		err = UpdateSignalSessionConfiguration(signalSessionConfiguration, customer)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// Registered signalSessionConfiguration
	newSignalSessionConfiguration = db_models.SignalSessionConfiguration{
		CustomerId:        customer.Id,
		SessionId:         signalSessionConfiguration.SessionId,
		HeartbeatInterval: signalSessionConfiguration.HeartbeatInterval,
		MaxRetransmit:     signalSessionConfiguration.MaxRetransmit,
		AckTimeout:        signalSessionConfiguration.AckTimeout,
		AckRandomFactor:   signalSessionConfiguration.AckRandomFactor,
	}
	_, err = session.Insert(&newSignalSessionConfiguration)
	if err != nil {
		log.Infof("signal_session_configuration insert err: %s", err)
		goto Rollback
	}

	// add Commit() after all actions
	err = session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

/*
 * Updates SignalSessionConfiguration in the DB.
 *
 * parameter:
 *  signalSessionConfiguration SignalSessionConfiguration
 *  customer Customer
 * return:
 *  err error
 */
func UpdateSignalSessionConfiguration(signalSessionConfiguration SignalSessionConfiguration, customer Customer) (err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	// Updated signalSessionConfiguration
	updSignalSessionConfiguration := new(db_models.SignalSessionConfiguration)
	_, err = engine.Where("customer_id = ? AND session_id = ?", customer.Id, signalSessionConfiguration.SessionId).Get(updSignalSessionConfiguration)
	if err != nil {
		return
	}
	if updSignalSessionConfiguration.Id == 0 {
		// no data found
		log.Infof("signal_session_configuration update data exitst err: %s", err)
		return
	}

	updSignalSessionConfiguration.HeartbeatInterval = signalSessionConfiguration.HeartbeatInterval
	updSignalSessionConfiguration.MaxRetransmit = signalSessionConfiguration.MaxRetransmit
	updSignalSessionConfiguration.AckTimeout = signalSessionConfiguration.AckTimeout
	updSignalSessionConfiguration.AckRandomFactor = signalSessionConfiguration.AckRandomFactor
	_, err = session.Where("id = ?", updSignalSessionConfiguration.Id).Update(updSignalSessionConfiguration)
	if err != nil {
		log.Infof("customer update err: %s", err)
		goto Rollback
	}

	// add Commit() after all actions
	err = session.Commit()
	return
Rollback:
	session.Rollback()
	return
}

/*
 * Obtains the SignalSessionConfiguration by the customer ID.
 *
 * parameter:
 *  customerId id of the customer
 *  sessionId id of the session
 * return:
 *  signalSessionConfiguration SignalSessionConfiguration
 *  error error
 */
func GetSignalSessionConfiguration(customerId int, sessionId int) (signalSessionConfiguration SignalSessionConfiguration, err error) {
	// default value setting
	signalSessionConfiguration = SignalSessionConfiguration{}

	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Printf("database connect error: %s", err)
		return
	}

	// Get customer table data
	dbSignalSessionConfiguration := db_models.SignalSessionConfiguration{}
	chk, err := engine.Where("customer_id = ? AND session_id = ?", customerId, sessionId).Get(&dbSignalSessionConfiguration)
	if err != nil {
		return
	}
	if !chk {
		// no data
		return
	}
	signalSessionConfiguration.SessionId = dbSignalSessionConfiguration.SessionId
	signalSessionConfiguration.HeartbeatInterval = dbSignalSessionConfiguration.HeartbeatInterval
	signalSessionConfiguration.MaxRetransmit = dbSignalSessionConfiguration.MaxRetransmit
	signalSessionConfiguration.AckTimeout = dbSignalSessionConfiguration.AckTimeout
	signalSessionConfiguration.AckRandomFactor = dbSignalSessionConfiguration.AckRandomFactor

	return

}

/*
 * Deletes the SignalSessionConfiguration by the customer ID and session id.
 *
 * parameter:
 *  customerId customer ID
 *  sessionId session ID
 * return:
 *  error error
 */
func DeleteSignalSessionConfiguration(customerId int, sessionId int) (err error) {
	// database connection create
	engine, err := ConnectDB()
	if err != nil {
		log.Errorf("database connect error: %s", err)
		return
	}

	// transaction start
	session := engine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return
	}

	dbSignalSessionConfiguration := db_models.SignalSessionConfiguration{}
	_, err = engine.Where("customer_id = ? AND session_id = ?", customerId, sessionId).Get(&dbSignalSessionConfiguration)
	if err != nil {
		log.Errorf("get signalSessionConfiguration err: %s", err)
		goto Rollback
	}

	// Delete signalSessionConfiguration table data
	_, err = session.Delete(db_models.SignalSessionConfiguration{CustomerId: customerId, SessionId: sessionId})
	if err != nil {
		log.Errorf("delete signalSessionConfiguration error: %s", err)
		goto Rollback
	}

	session.Commit()
	return
Rollback:
	session.Rollback()
	return
}
