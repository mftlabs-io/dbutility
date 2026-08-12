package dbutility

import (
	"amfui/dbconnector"
	"amfui/utilities"
	"time"
)

type DbUtil struct {
	Db       *dbconnector.DbConnector
	Timezone string
}

func (util *DbUtil) PrepareQuery(context utilities.AppContext, Db *dbconnector.DbConnector, utilitytype, startdate, enddate string, clean, validatemain, validatehistory bool) error {
	startdate = startdate + " 00:00:00.000"
	enddate = enddate + " 23:59:59.999"
	now := time.Now().UTC()
	presentday := now.AddDate(0, 0, -1)
	previousdays := presentday.AddDate(0, 0, -14)
	fromdate := previousdays.Format("2006-01-02") + " 23:59:59.999"
	if clean == false && utilitytype == "all" && validatemain == false && validatehistory == false {
		context.Logger.Info("Date range is::%v\n", fromdate)
		// RangeAll now moves rows (archive + delete) in a single atomic
		// statement per table, so no separate DeleteAll pass is needed.
		err := util.RangeAll(context, Db, fromdate)
		if err != nil {
			context.Logger.Info("error when moving records to history tables%v", err)
			return err
		}
	} else if utilitytype == "" && clean == false && validatemain == false && validatehistory == false {
		context.Logger.Info("Start Date is::%v\n", startdate)
		context.Logger.Info("End Date is::%v\n", enddate)
		// WithinRange now moves rows (archive + delete) in a single atomic
		// statement per table, so no separate DeleteWithinRange pass is needed.
		err := util.WithinRange(context, Db, startdate, enddate)
		if err != nil {
			context.Logger.Info("error when moving records to history tables%v\n", err)
			return err
		}
	} else if clean && utilitytype == "all" && validatemain == false && validatehistory == false {
		context.Logger.Info("Date range is::%v\n", fromdate)
		err := util.DeleteAll(context, Db, fromdate)
		if err != nil {
			context.Logger.Info("error when inserting range of records to history tables%v\n", err)
			return err
		}
	} else if clean && utilitytype == "" && validatemain == false && validatehistory == false {
		context.Logger.Info("Start Date is::%v\n", startdate)
		context.Logger.Info("End Date is::%v\n", enddate)
		err := util.DeleteWithinRange(context, Db, startdate, enddate)
		if err != nil {
			context.Logger.Info("error when inserting range of records to history tables%v\n", err)
			return err
		}
	} else if validatemain && utilitytype == "all" {
		err := util.ValidateAll(context, Db, fromdate, "main")
		if err != nil {
			// SEC-010: silent failure left the CLI looking like a no-op; propagate so the operator sees it
			context.Logger.Warn("PrepareQuery: ValidateAll(main, fromdate=%v) failed: %v", fromdate, err)
			return err
		}
	} else if validatemain && utilitytype == "" {
		err := util.ValidateWithinRange(context, Db, startdate, enddate, "main")
		if err != nil {
			// SEC-010: silent failure left the CLI looking like a no-op; propagate so the operator sees it
			context.Logger.Warn("PrepareQuery: ValidateWithinRange(main, %v..%v) failed: %v", startdate, enddate, err)
			return err
		}
	} else if validatehistory && utilitytype == "all" {
		err := util.ValidateAll(context, Db, fromdate, "history")
		if err != nil {
			// SEC-010: silent failure left the CLI looking like a no-op; propagate so the operator sees it
			context.Logger.Warn("PrepareQuery: ValidateAll(history, fromdate=%v) failed: %v", fromdate, err)
			return err
		}
	} else if validatehistory && utilitytype == "" {
		err := util.ValidateWithinRange(context, Db, startdate, enddate, "history")
		if err != nil {
			// SEC-010: silent failure left the CLI looking like a no-op; propagate so the operator sees it
			context.Logger.Warn("PrepareQuery: ValidateWithinRange(history, %v..%v) failed: %v", startdate, enddate, err)
			return err
		}
	}

	return nil
}

// RangeAll moves all rows older than last14daydate from each main table to its
// history table. Each table is one atomic delete+insert statement (MoveToHistory),
// which replaces the old count / insert-select / sleep / delete sequence.
func (util *DbUtil) RangeAll(context utilities.AppContext, Db *dbconnector.DbConnector, last14daydate string) error {
	context.Logger.Info("Date range is: %v\n", last14daydate)
	for _, table := range archiveOrder {
		err := util.MoveToHistory(context, Db, "", last14daydate, table)
		if err != nil {
			context.Logger.Info("error when moving %v records to history table%v", table, err)
			return err
		}
	}
	return nil
}

// WithinRange moves all rows between startdate and enddate from each main table
// to its history table, one atomic delete+insert statement per table.
func (util *DbUtil) WithinRange(context utilities.AppContext, Db *dbconnector.DbConnector, startdate, enddate string) error {
	for _, table := range archiveOrder {
		err := util.MoveToHistory(context, Db, startdate, enddate, table)
		if err != nil {
			context.Logger.Info("error when moving %v records to history table%v", table, err)
			return err
		}
	}
	return nil
}

// DeleteAll removes rows older than last14daydate from the main tables without
// archiving them (clean-only mode). DeleteHistory logs rows affected, so the
// former per-table count(*) pre-scans were dropped.
func (util *DbUtil) DeleteAll(context utilities.AppContext, Db *dbconnector.DbConnector, last14daydate string) error {
	for _, table := range archiveOrder {
		err := util.DeleteHistory(context, Db, "", last14daydate, table)
		if err != nil {
			context.Logger.Info("error when deleting %v table%v", table, err)
			return err
		}
	}
	return nil
}

// DeleteWithinRange removes rows between the two dates from the main tables
// without archiving them (clean-only mode).
func (util *DbUtil) DeleteWithinRange(context utilities.AppContext, Db *dbconnector.DbConnector, last14daydate, presentDate string) error {
	for _, table := range archiveOrder {
		err := util.DeleteHistory(context, Db, last14daydate, presentDate, table)
		if err != nil {
			context.Logger.Info("error when deleting %v table%v", table, err)
			return err
		}
	}
	return nil
}

func (util *DbUtil) ValidateAll(context utilities.AppContext, Db *dbconnector.DbConnector, last14daydate, tabletype string) error {
	var tablename string
	var sessiontable string
	var sessionreltable string
	var eventtable string
	if tabletype == "main" {
		tablename = "amf_message"
		sessiontable = "amf_session"
		sessionreltable = "amf_session_rel"
		eventtable = "amf_event"
	} else {
		tablename = "amf_message_history"
		sessiontable = "amf_session_history"
		sessionreltable = "amf_session_rel_history"
		eventtable = "amf_event_history"
	}
	count, cmerr := util.CheckCount(context, Db, tablename, "", last14daydate)
	if cmerr != nil {
		// SEC-010: ValidateAll counts are informational; log + continue
		context.Logger.Warn("ValidateAll: CheckCount(%v) failed: %v", tablename, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", tablename, count)
	count1, cmerr := util.CheckCount(context, Db, sessiontable, "", last14daydate)
	if cmerr != nil {
		context.Logger.Warn("ValidateAll: CheckCount(%v) failed: %v", sessiontable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", sessiontable, count1)
	count2, cmerr := util.CheckCount(context, Db, sessionreltable, "", last14daydate)
	if cmerr != nil {
		context.Logger.Warn("ValidateAll: CheckCount(%v) failed: %v", sessionreltable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", sessionreltable, count2)
	count3, cmerr := util.CheckCount(context, Db, eventtable, "", last14daydate)
	if cmerr != nil {
		context.Logger.Warn("ValidateAll: CheckCount(%v) failed: %v", eventtable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", eventtable, count3)
	util.CheckDistinctSenderWithCount(context, Db, tablename, "", last14daydate)
	util.CheckDistinctReceiverWithCount(context, Db, tablename, "", last14daydate)
	return nil
}

func (util *DbUtil) ValidateWithinRange(context utilities.AppContext, Db *dbconnector.DbConnector, startdate, enddate, tabletype string) error {
	var tablename string
	var sessiontable string
	var sessionreltable string
	var eventtable string
	if tabletype == "main" {
		tablename = "amf_message"
		sessiontable = "amf_session"
		sessionreltable = "amf_session_rel"
		eventtable = "amf_event"
	} else {
		tablename = "amf_message_history"
		sessiontable = "amf_session_history"
		sessionreltable = "amf_session_rel_history"
		eventtable = "amf_event_history"
	}
	count, cmerr := util.CheckCount(context, Db, tablename, startdate, enddate)
	if cmerr != nil {
		// SEC-010: ValidateWithinRange counts are informational; log + continue
		context.Logger.Warn("ValidateWithinRange: CheckCount(%v) failed: %v", tablename, cmerr)
	}
	context.Logger.Info("Count from message table is: %v\n", count)
	count1, cmerr := util.CheckCount(context, Db, sessiontable, startdate, enddate)
	if cmerr != nil {
		context.Logger.Warn("ValidateWithinRange: CheckCount(%v) failed: %v", sessiontable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", sessiontable, count1)
	count2, cmerr := util.CheckCount(context, Db, sessionreltable, startdate, enddate)
	if cmerr != nil {
		context.Logger.Warn("ValidateWithinRange: CheckCount(%v) failed: %v", sessionreltable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", sessionreltable, count2)
	count3, cmerr := util.CheckCount(context, Db, eventtable, startdate, enddate)
	if cmerr != nil {
		context.Logger.Warn("ValidateWithinRange: CheckCount(%v) failed: %v", eventtable, cmerr)
	}
	context.Logger.Info("Count from %v table is: %v\n", eventtable, count3)
	util.CheckDistinctSenderWithCount(context, Db, tablename, startdate, enddate)
	util.CheckDistinctReceiverWithCount(context, Db, tablename, startdate, enddate)
	return nil
}
