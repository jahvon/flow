import {
  Anchor,
  Button,
  Container,
  Group,
  Image,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { IconBook, IconBrandGithub } from "@tabler/icons-react";
import iconImage from "../../assets/icon.png";
import styles from "./Welcome.module.css";

interface WelcomeProps {
  welcomeMessage?: string;
}

export function Welcome({ welcomeMessage }: WelcomeProps) {
  return (
    <Container size="xs" className={styles.welcome}>
      <Stack align="center" gap={25}>
        <Stack align="center" gap={10}>
          <Image
            src={iconImage}
            alt="flow"
            width={120}
            height={120}
            fit="contain"
            m={60}
            className={styles.welcome__logo}
          />
          <Title order={3} fw={300} className={styles.welcome__title}>
            flow desktop
          </Title>
        </Stack>

        <Group gap={20}>
          <Anchor href="https://flowexec.io" target="_blank" underline="never">
            <Button
              size="xs"
              variant="outline"
              leftSection={<IconBook size={12} />}
              className={styles.welcome__button}
            >
              Docs
            </Button>
          </Anchor>
          <Anchor
            href="https://github.com/jahvon/flow"
            target="_blank"
            underline="never"
          >
            <Button
              size="xs"
              variant="outline"
              leftSection={<IconBrandGithub size={12} />}
              className={styles.welcome__button}
            >
              GitHub
            </Button>
          </Anchor>
        </Group>

        {welcomeMessage && (
          <Text size="md" c="dimmed" ta="center" m={10}>
            {welcomeMessage}
          </Text>
        )}
      </Stack>
    </Container>
  );
}
