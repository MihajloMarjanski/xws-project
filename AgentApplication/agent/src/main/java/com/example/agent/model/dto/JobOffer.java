package com.example.agent.model.dto;

import com.example.agent.model.JobPosition;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class JobOffer {
    private String jobPosition;
    private String companyName;
    private String jobInfo;
    private String qualifications;
}
