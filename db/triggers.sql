--Triggerid peab ikka käsitsi sisestama terminali
CREATE TRIGGER new_comment
AFTER INSERT ON comment
BEGIN
    INSERT INTO notification(user_id, comment_id)
    VALUES(new.user_id, new.id);
END;


CREATE TRIGGER new_vote
AFTER INSERT ON vote
BEGIN
    INSERT INTO notification(user_id, vote_id)
    VALUES(new.user_id, new.id);
END;

CREATE TRIGGER new_report
AFTER INSERT ON report
BEGIN
    INSERT INTO notification(user_id, report_id)
    VALUES(new.user_id, new.reporter_id);
END;

CREATE TRIGGER new_upgrade
AFTER INSERT ON upgrade_request
BEGIN
    INSERT INTO notification(user_id, upgrade_id)
    VALUES(new.user_id, new.requester_id);
END;