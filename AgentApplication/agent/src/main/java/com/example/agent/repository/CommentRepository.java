package com.example.agent.repository;

import com.example.agent.model.Comment;
import com.example.agent.model.JobPosition;
import org.springframework.data.jpa.repository.JpaRepository;

public interface CommentRepository extends JpaRepository<Comment, Integer> {
}
