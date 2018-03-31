SET citus.task_executor_type TO "task-tracker";

SELECT sum(systemevents.timestamp - usercommand.timestamp), usercommand.command from usercommand 
join systemevents on usercommand.transactionnum = systemevents.transactionnum GROUP BY usercommand.command ORDER BY sum DESC;

select timestamp - quoteservertime - 248776 from quoteserver as ms
docker exec  5fd65650e62c psql -U postgres -d logs -t -A -F"," -c "select timestamp - quoteservertime - 248776 from quoteserver as ms;" > query.csv
docker exec  5fd65650e62c psql -U postgres -d logs -t -A -F"," -c "SET citus.task_executor_type TO 'task-tracker'; SELECT systemevents.timestamp - usercommand.timestamp from usercommand join systemevents on usercommand.transactionnum = systemevents.transactionnum WHERE usercommand.command = 'QUOTE';" > quote_times.csv
