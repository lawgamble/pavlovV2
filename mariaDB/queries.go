package mariadb

var tempTeamQuery = `
SELECT
    T.TeamName AS Team,
    T.TeamRegion AS Region,
    GROUP_CONCAT(U.DiscordName ORDER BY U.DiscordName) AS PlayerNames,
    GROUP_CONCAT(U.PlayStyle ORDER BY U.DiscordName) AS PlayStyles,
    GROUP_CONCAT(U.Region ORDER BY U.DiscordName) AS PlayerRegions,
    GROUP_CONCAT(U.InGameName ORDER BY U.DiscordName) AS InGameNames,
    GROUP_CONCAT(U.PlayerType ORDER BY U.DiscordName) AS PlayerTypes,
    GROUP_CONCAT(U.DOB ORDER BY U.DiscordName) AS DOBs
FROM
    SND_TEAMS T
JOIN
    SND_TEMP_ROSTERS TR ON T.TeamName = TR.Team
JOIN
    SND_USERS U ON TR.DiscordId = U.DiscordId
WHERE
    T.TeamStatus = 'Pending'
GROUP BY
    T.TeamName;

`
