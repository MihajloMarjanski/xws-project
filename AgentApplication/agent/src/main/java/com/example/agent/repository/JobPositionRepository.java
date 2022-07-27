package com.example.agent.repository;

import com.example.agent.model.JobPosition;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface JobPositionRepository extends JpaRepository<JobPosition, Integer> {
    List<JobPosition> findAllByCompanyId(Integer companyId);
}
