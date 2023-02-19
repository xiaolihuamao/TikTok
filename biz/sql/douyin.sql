# 建立数据库

# create database douyin;
# use douyin;

# 建立模式

# create schema douyin;
 use douyin;


drop table if exists  users;
create table users(
                      user_id bigint auto_increment primary key , #用户id（主键上自动建立索引）
                      username varchar(32)  character set utf8mb4 collate utf8mb4_da_0900_ai_ci not null unique ,   #用户名（注册中要求用户名唯一）
                      password varchar(300) character set utf8mb4 collate utf8mb4_da_0900_ai_ci not null ,   #密码(可能有md5或其他加密操作，字符串比较长）
                      follow_count bigint default 0,  #关注总数
                      follower_count bigint default 0,  #粉丝总数
                      unique index index_name(username) using btree   #在用户名上创建唯一索引（注册中要求用户名唯一）
)engine =InnoDB auto_increment=0 character set utf8mb4 collate utf8mb4_da_0900_ai_ci row_format =dynamic
    comment '用户表'
;

drop table if exists videos;
create table videos(
                       video_id bigint auto_increment primary key ,  #视频id
                       author_id bigint not null ,   #视频作者id
                       title varchar(300) character set utf8mb4 collate utf8mb4_da_0900_ai_ci null default "", #视频标题
                       play_url varchar(500) character set utf8mb4 collate utf8mb4_da_0900_ai_ci null default "",  #视频播放地址
                       cover_url varchar(500) character set utf8mb4 collate utf8mb4_da_0900_ai_ci null default "",  #视频封面地址
                       create_time datetime not null default current_timestamp, #视频的最新投稿时间戳，精确到秒，不填表示当前时间
                       favorite_count bigint not null default 0, #视频点赞人数
                       comment_count bigint not null default 0,  #视频评论人数
                       is_favorite tinyint not null default 0, #是否点赞（mysql数据库的布尔型对应的数据类型为tinyint(1),存入的数据,0代表false（未点赞）,1代表true （已点赞）
                       index index_id(author_id) using btree ,  #在外键上创建索引
                       foreign key (author_id) references users(user_id) on delete cascade on update restrict #用户id作为外键(级联操作，删除视频表中对应的数据，不存在外键删除问题）
)engine =InnoDB auto_increment=0 character set utf8mb4 collate utf8mb4_da_0900_ai_ci row_format =dynamic
    comment '视频表'
;

drop table if exists comments;
create table comments (
                          comment_id bigint auto_increment primary key, #评论id
                          video_id bigint not null , #视频id
                          user_id bigint not null ,  #用户id
                          content varchar(200) character set utf8mb4 collate utf8mb4_da_0900_ai_ci null default "",  #评论内容
                          create_date varchar(10) character set utf8mb4 collate utf8mb4_da_0900_ai_ci null default "", #评论发布日期（格式 mm-dd)
                          index index_vid(video_id) using btree ,  #在外键上创建索引
                          index index_uid(user_id) using btree ,  #在外键上创建索引
                          foreign key (video_id) references videos(video_id) on delete cascade on update restrict,  #视频id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
                          foreign key (user_id) references users(user_id) on delete cascade on update restrict  #用户id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
)engine =InnoDB auto_increment=0 character set utf8mb4 collate utf8mb4_da_0900_ai_ci row_format =dynamic
    comment '评论表'
;

drop table if exists favorites;
create table favorites(
                          favorite_id bigint auto_increment primary key ,  #点赞id
                          video_id bigint , #视频id
                          user_id bigint ,  #用户id
                          index index_vid(video_id) using btree ,  #在外键上创建索引
                          index index_uid(user_id) using btree ,  #在外键上创建索引
                          foreign key (video_id) references videos(video_id) on delete cascade on update restrict,  #视频id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
                          foreign key (user_id) references users(user_id) on delete cascade on update restrict #用户id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
)engine =InnoDB auto_increment=0 character set utf8mb4 collate utf8mb4_da_0900_ai_ci row_format =dynamic
    comment '点赞表'
;

drop table if exists relations;
create table relations(
                          relation_id bigint auto_increment primary key ,  #关注id
                          follow_id bigint not null , #用户id
                          follower_id bigint not null ,  #用户关注者id（别人关注我）
                          index index_uid(follow_id) using btree ,  #在外键上创建索引
                          index index_uid1(follower_id) using btree ,  #在外键上创建索引
                          foreign key (follow_id) references users(user_id) on delete cascade on update restrict, #用户id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
                          foreign key (follower_id) references users(user_id) on delete cascade on update restrict #用户id作为外键(级联操作，删除评论表中对应的数据，不存在外键删除问题）
)engine =InnoDB auto_increment=0 character set utf8mb4 collate utf8mb4_da_0900_ai_ci row_format =dynamic
    comment '关注表'
;
