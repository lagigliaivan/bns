package ar.com.bestprice.buyitnow;

import org.junit.Test;

import java.util.Calendar;
import java.util.Date;

import static org.junit.Assert.*;

/**
 * To work on unit tests, switch the Test Artifact in the Build Variants view.
 */
public class ExampleUnitTest {
    @Test
    public void addition_isCorrect() throws Exception {


        Calendar purchaseDateTime = Calendar.getInstance();
        StringBuffer strigBuffer = new StringBuffer();

        Date date = new Date(Long.parseLong("1469992368") * 1000);
        purchaseDateTime.setTime(date);

        strigBuffer.append("Dia: ");
        strigBuffer.append(purchaseDateTime.get(Calendar.DAY_OF_MONTH));
        strigBuffer.append("/");
        strigBuffer.append(purchaseDateTime.get(Calendar.MONTH));

        strigBuffer.append(" ");
        strigBuffer.append(purchaseDateTime.get(Calendar.HOUR));
        strigBuffer.append(":");
        strigBuffer.append(purchaseDateTime.get(Calendar.MINUTE));

        System.out.println(strigBuffer.toString());

    }

}
