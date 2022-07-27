package com.example.agent.repository;

import com.example.agent.model.InterviewInformation;
import com.example.agent.model.JobPosition;
import org.springframework.data.jpa.repository.JpaRepository;

public interface InterviewInformationRepository extends JpaRepository<InterviewInformation, Integer> {
}
