SET citus.task_executor_type TO "task-tracker";
SELECT sum(systemevents.timestamp - usercommand.timestamp), usercommand.command from usercommand 
join systemevents on usercommand.transactionnum = systemevents.transactionnum GROUP BY usercommand.command ORDER BY sum DESC;