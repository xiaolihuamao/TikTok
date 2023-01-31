drop procedure if exists addRoomPrice;
create procedure addRoomPrice()
begin
    declare i int default 0;
    set i = 0;
    start transaction
        ;
        while i < 100 do
                -- insert sql
                insert into videos(video_id, author_id, title, create_time)
                values (default, 1, '视频',date_sub('2021-01-01', interval -i day));

                set i = i + 1;
            end while;
    commit;
end
