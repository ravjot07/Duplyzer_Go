// src/pages/index.tsx
import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {
  Container,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Input,
  Heading,
  Spinner,
  Box,
} from '@chakra-ui/react';

interface ReportEntry {
  hash: string;
  files: string[];
}

const Home: React.FC = () => {
  const [report, setReport] = useState<ReportEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    axios.get('http://localhost:8080/report')
      .then(response => {
        setReport(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching report data:', error);
        setLoading(false);
      });
  }, []);

  const filteredReport = report.filter(entry =>
    entry.files.some(file => file.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  if (loading) {
    return <Spinner size="xl" />;
  }

  return (
    <Container maxW="container.xl" py={8}>
      <Heading as="h1" size="xl" mb={8} textAlign="center">
        Duplicate File Report
      </Heading>
      <Input
        placeholder="Search files..."
        value={searchTerm}
        onChange={e => setSearchTerm(e.target.value)}
        mb={4}
      />
      <Box overflowX="auto">
        <Table variant="simple">
          <Thead>
            <Tr>
              <Th>Hash</Th>
              <Th>Files</Th>
            </Tr>
          </Thead>
          <Tbody>
            {filteredReport.map((entry, index) => (
              <Tr key={index}>
                <Td>{entry.hash}</Td>
                <Td>
                  {entry.files.map((file, fileIndex) => (
                    <Box key={fileIndex}>{file}</Box>
                  ))}
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </Box>
    </Container>
  );
};

export default Home;
