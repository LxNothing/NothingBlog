CREATE DATABASE IF NOT EXISTS ngblog CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
use ngblog;

CREATE TABLE IF NOT EXISTS `user`(
    `id` bigint(20) not null auto_increment comment '自增id',
    `user_id` bigint(20) unique not null comment '应用层生成的用户id',
    `username` varchar(64) collate utf8mb4_general_ci not null,
    `password` varchar(64) collate utf8mb4_general_ci not null,
    `email` varchar(64) unique collate utf8mb4_general_ci,
    `gender` tinyint(4) not null default 0,
    `create_time` timestamp null default current_timestamp comment '创建时间',
    `update_time` timestamp null default current_timestamp on update current_timestamp comment '更新时间',
    primary key (`id`),
    unique key `idx_username` (`username`) using btree,
    unique key `idx_user_id` (`user_id`) using btree
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `class` (
    `id` int not null auto_increment,
    `name` varchar(64) not null comment 'class name',
    `article_count` int not null default 0 comment '该标签下的文章数量',
    `desc` varchar(256) comment '标签的描述信息',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `article` (
    `id` bigint(20) not null auto_increment,
    `article_id` bigint(20) not null comment '应用层生成的文章的id',
    `title` varchar(256) collate utf8mb4_general_ci not null comment '文章标题',
    `content` MEDIUMTEXT collate utf8mb4_general_ci not null comment '文章内容',
    `author_id` bigint(20) not null comment '作者id',
    `class_id` bigint(20) not null comment '文章所属的类别id',
    `status` tinyint(4) not null default 1 comment '博客状态',
    `comment_count` int not null default 0 comment '评论数量',
    `visit_count` int not null default 0 comment '投票分数',
    `vote_count` int not null default 0 comment '点赞数量',
    `del_falg` char not null default 0 comment '删除标记 0-未删除 1-删除',
    `en_comment` tinyint not null default 0 comment '是否允许评论 0 - 不允许 1-允许',
    `authority` tinyint not null default 0 comment '文章权限 0 - 公开 1 - 私有 2 - 指定用户可见(2先不做)',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`),
    unique key `idx_article_id` (`article_id`) using btree,
    key `idx_author_id` (`author_id`) using btree,
    constraint foreign key (author_id) references user(user_id),
    constraint foreign key (class_id) references user(id)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tag` (
    `id` int not null auto_increment,
    `name` varchar(64) not null comment 'tag name',
    `article_count` int not null default 0 comment '该标签下的文章数量',
    `desc` varchar(256) comment '标签的描述信息',
    `create_time` timestamp null default current_timestamp,
    `update_time` timestamp null default current_timestamp on update current_timestamp,
    primary key (`id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `tag_article` (
    `id` int not null auto_increment,
    `tag_id` int comment 'tag id',
    `article_id` bigint(20) comment '文章id',
    constraint foreign key (tag_id) references tag(id),
    constraint foreign key (article_id) references article(article_id),
    primary key (`id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `comment` (
    `id` int not null auto_increment,
    `article_id` bigint(20) comment '所属的文章id',
    `comment_id` bigint(20) not null comment '应用层生成的评论id',
    `create_uid` bigint(20) not null comment '创建该条评论的用户id',
    `root_comm_id` bigint(20) not null default -1 comment '所属根评论ID',
    `parent_comm_uid` bigint(20) not null comment '所属的父评论的用户ID',
    `parent_comm_id` bigint(20) not null comment '所属的父评论的评论ID',
    `content` text comment '评论内容',
    `type` tinyint not null comment '评论类型 0-文章评论 1-其他评论',
    `del_flag` tinyint not null default 0 comment '删除标记 0-未删除 1-删除',
    constraint foreign key (article_id) references article(article_id),
    constraint foreign key (create_uid) references user(user_id),
    primary key (`id`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;
