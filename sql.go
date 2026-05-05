package dbutility

import (
	"amfui/utilities"
	"amfui/dbconnector"
)

const (
	SELECT_MESSAGES = "select sender,receiver,msg_type,control_no,file_name,file_type,file_path,workflow_id,session_id,parent_id,doc_count,origin,reference_id,status,status_time::timestamptz,can_requeue,can_reprocess,create_time::timestamptz,created_by,message_id,file_size,site_id,node_id,can_req_and_rep,data_type,contents from amf_message msg "
	INSERT_INTO_MESSAGE_HISTORY = "insert into amf_message_history (sender,receiver,msg_type,control_no,file_name,file_type,file_path,workflow_id,session_id,parent_id,doc_count,origin,reference_id,status,status_time,can_requeue,can_reprocess,create_time,created_by,message_id,file_size,site_id,node_id,can_req_and_rep,data_type,contents) "

	SELECT_SESSIONS = "select session_id,session_start,session_end,workflow_name,instance_id,username,status,create_time,created_by,site_id,node_id from amf_session"
	INSERT_INTO_SESSION_HISTORY = "insert into amf_session_history (session_id,session_start,session_end,workflow_name,instance_id,username,status,create_time,created_by,site_id,node_id)"

	SELECT_SESSION_REL = "select relation_id,session_id,message_id,rel_type,create_time,created_by from amf_session_rel"
	INSERT_INTO_SESSION_REL_HISTORY = "insert into amf_session_rel_history (relation_id,session_id,message_id,rel_type,create_time,created_by)"

	SELECT_EVENT = "select event_id,level,message_id,session_id,action_id,text,status,create_time,created_by from amf_event"
	INSERT_INTO_EVENT_HISTORY = "insert into amf_event_history (event_id,level,message_id,session_id,action_id,text,status,create_time,created_by)"
)

func (util *DbUtil) DeleteHistory(context utilities.AppContext,Db *dbconnector.DbConnector,fromdate,todate,tablename string) error{
	if fromdate == ""{
		query := "delete from "+tablename+" where create_time <= $1"
		resp, err := Db.Exec(query,todate)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		//context.Logger.Info("Delete query for %v is: %v\n",tablename,query)
		context.Logger.Info("Delete response for %v is: %v\n",tablename,rowcount)
	} else {
		query := "delete from "+tablename+" where create_time >= $1 and create_time <= $2"
		resp, err := Db.Exec(query,fromdate,todate)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		//context.Logger.Info("Delete query for %v is: %v\n",tablename,query)
		context.Logger.Info("Delete response for %v is: %v\n",tablename,rowcount)
	}

	return nil
}

func (util *DbUtil) InsertToHistoryTable(context utilities.AppContext,Db *dbconnector.DbConnector,then string,tablename string) error{
	if tablename == "amf_message_history"{
		Query := INSERT_INTO_MESSAGE_HISTORY+SELECT_MESSAGES+" where create_time <= $1"
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query,then)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_session_history"{
		Query := INSERT_INTO_SESSION_HISTORY+SELECT_SESSIONS+" where create_time <= $1"
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query,then)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_session_rel_history"{
		Query := INSERT_INTO_SESSION_REL_HISTORY+SELECT_SESSION_REL+" where create_time <= $1"
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query,then)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_event_history"{
		Query := INSERT_INTO_EVENT_HISTORY+SELECT_EVENT+" where create_time <= $1"
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query,then)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	}
	return nil
}

func (util *DbUtil) InsertLastMonthHistory(context utilities.AppContext,Db *dbconnector.DbConnector,last14daydate,presentDate,tablename string) error{
	if tablename == "amf_message_history"{
		whereClause := " where create_time >= '"+last14daydate+"' and create_time <= '"+presentDate+"'"
		Query := INSERT_INTO_MESSAGE_HISTORY+SELECT_MESSAGES+whereClause
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_session_history"{
		whereClause := " where create_time >= '"+last14daydate+"' and create_time <= '"+presentDate+"'"
		Query := INSERT_INTO_SESSION_HISTORY+SELECT_SESSIONS+whereClause
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_session_rel_history"{
		whereClause := " where create_time >= '"+last14daydate+"' and create_time <= '"+presentDate+"'"
		Query := INSERT_INTO_SESSION_REL_HISTORY+SELECT_SESSION_REL+whereClause
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	} else if tablename == "amf_event_history"{
		whereClause := " where create_time >= '"+last14daydate+"' and create_time <= '"+presentDate+"'"
		Query := INSERT_INTO_EVENT_HISTORY+SELECT_EVENT+whereClause
		context.Logger.Info("Insert query for %v is: %v\n",tablename,Query)
		resp, err := Db.Exec(Query)
		if err != nil{
			return err
		}
		rowcount,_ := resp.RowsAffected()
		context.Logger.Info("response for adding data to %v is: %v\n",tablename,rowcount)
	}
	return nil
}

func (util *DbUtil) CheckCount(context utilities.AppContext,Db *dbconnector.DbConnector,tablename,fromdate,todate string) (int,error){
	if fromdate == ""{
		query := "select count(*) from "+tablename+" where create_time <= $1"
		row := Db.QueryRow(query,todate)
		var count int
		err := row.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	} else {
		query := "select count(*) from "+tablename+" where create_time >= $1 and create_time <= $2"
		row := Db.QueryRow(query,fromdate,todate)
		var count int
		err := row.Scan(&count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0,nil

}

func (util *DbUtil) CheckDistinctSenderWithCount (context utilities.AppContext,Db *dbconnector.DbConnector,tablename,fromdate,todate string){
	if fromdate == ""{
		query := "select distinct sender,count(*) from "+tablename+" where create_time <= $1 group by sender"
		rows,err := Db.Query(query,todate)
		if err != nil{
			//return nil,err
		}
		defer rows.Close()
		for rows.Next() {
			var count int
			var sender string
			err := rows.Scan(&sender,&count)
			if err != nil {
				//return 0, err
			}
			context.Logger.Info("Sender :%v with Count: %v from %v\n",sender,count,tablename)
			context.Logger.Info("CheckDistinctSenderWithCount: sender=%v count=%v table=%v", sender, count, tablename)

		}
	} else {
		query := "select distinct sender,count(*) from "+tablename+" where create_time >= $1 and create_time <= $2 group by sender"
		rows,err := Db.Query(query,fromdate,todate)
		if err != nil{
			//return nil,err
		}
		defer rows.Close()
		for rows.Next() {
			var count int
			var sender string
			err := rows.Scan(&sender,&count)
			if err != nil {
				//return 0, err
			}
			context.Logger.Info("Sender :%v with Count: %v from %v\n",sender,count,tablename)
			context.Logger.Info("CheckDistinctSenderWithCount: sender=%v count=%v table=%v", sender, count, tablename)
		}
	}
}

func (util *DbUtil) CheckDistinctReceiverWithCount (context utilities.AppContext,Db *dbconnector.DbConnector,tablename,fromdate,todate string){
	if fromdate == ""{
		query := "select distinct receiver,count(*) from "+tablename+" where create_time <= $1 group by receiver"
		rows,err := Db.Query(query,todate)
		if err != nil{
			//return nil,err
		}
		defer rows.Close()
		for rows.Next() {
			var count int
			var receiver string
			err := rows.Scan(&receiver,&count)
			if err != nil {
				//return 0, err
			}
			context.Logger.Info("Receiver :%v with Count: %v from %v\n",receiver,count,tablename)
			context.Logger.Info("CheckDistinctReceiverWithCount: receiver=%v count=%v table=%v", receiver, count, tablename)

		}
	} else {
		query := "select distinct receiver,count(*) from "+tablename+" where create_time >= $1 and create_time <= $2 group by receiver"
		rows,err := Db.Query(query,fromdate,todate)
		if err != nil{
			//return nil,err
		}
		defer rows.Close()
		for rows.Next() {
			var count int
			var receiver string
			err := rows.Scan(&receiver,&count)
			if err != nil {
				//return 0, err
			}
			context.Logger.Info("Receiver :%v with Count: %v from %v\n",receiver,count,tablename)
			context.Logger.Info("CheckDistinctReceiverWithCount: receiver=%v count=%v table=%v", receiver, count, tablename)
		}
	}
}