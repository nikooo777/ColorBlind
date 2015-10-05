package application;

import java.util.Random;

import javafx.application.Application;
import javafx.stage.Stage;

public class Main extends Application
{

	static Random rand;

	@Override
	public void start(final Stage primaryStage /* which is the main window */)
	{
		final ColorBlindWindow cbw = new ColorBlindWindow();
		cbw.show();
	}

	// main method
	public static void main(final String[] args)
	{
		// initialize the random
		rand = new Random();
		rand.setSeed(System.currentTimeMillis());

		// launch javaFX
		launch(args);
	}
}
